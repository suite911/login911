package login

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"path"
	"strconv"

	"github.com/suite911/error911"

	"github.com/pkg/browser"
	pkgErrors "github.com/pkg/errors"
	"github.com/valyala/fasthttp"
	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/sha3"
)

type Account struct {
	ID          int64  `json:"id"`
	RowID       int64  `json:"rowid"`
	Salt        []byte `json:"salt"`
	HasPassword bool   `json:"has_pw"`
}

type Token struct {
	ID         int64     `json:"id"`
	RowID      int64     `json:"rowid"`
	Email      string    `json:"email"`
	Username   string    `json:"username"`
	SessionKey []byte    `json:"session_key"`
	Expiration time.Time `json:"expiration"`
}

func LogIn(host string) (*Token, error) {
	var email, username, msg string
	var button int8
	user := path.Join(host, "api", "user")
	logIn := path.Join(host, "api", "login")
	for {
		dialog.New("email", &email, &username, &button, msg)
		switch button {
		case dialog.Cancel:
			return nil, error911.NewCancel()
		case dialog.Register:
			if err := register(host, email, username); err != nil {
				return nil, err
			}
			continue
		}
		var args fasthttp.Args
		args.Set("email", email)
		args.Set("username", username)
		statusCode, body, err := fasthttp.Post(nil, user, &args)
		if err != nil {
			return nil, err
		}
		switch statusCode {
		case 200:
		case 404:
			msg = "E-mail or username not found." // TODO: translate
			continue
		default:
			return nil, pkgErrors.Wrap(errors.New("HTTP "+strconv.Itoa(statusCode)), user)
		}

		var account Account
		if err := json.Unmarshal(body, &account); err != nil {
			return nil, err
		}
		var token *Token
		for {
			var password, retype string
			if account.HasPassword {
				// old account -- ask for their password
				dialog.New("password", &password, &button)
			} else {
				// new account -- ask them to create a password
				dialog.New("password", &password, &retype, &button)
			}
			if button == dialog.Cancel {
				break
			}
			if account.HasPassword {
				// old account -- verify their password
			} else {
				// new account -- set their password
				if retype != password {
					msg = "The passwords don't match." // TODO: translate
					continue
				}
			}
			key := argon2.IDKey([]byte(password), account.Salt, 1, 64*1024, 4, 32)
			var args fasthttp.Args
			args.Set("rowid", strconv.FormatInt(account.RowID, 10))
			args.Set("id", strconv.FormatInt(account.ID, 10))
			if account.HasPassword {
				// old account -- verify their password
				//args.Set("key", key)
				buf := make([]byte, 32, 32 + len(key))
				if _, err := rand.Read(rnd); err != nil {
					return err
				}
				hexRand := make([]byte, hex.EncodedLen(32))
				hex.Encode(hexRand, buf[:32])
				args.SetBytesV("rand", string(hexRand))
				buf = append(buf, key...)
				dig := sha3.Sum256(buf)
				hexDig := make([]byte, hex.EncodedLen(len(dig)))
				hex.Encode(hexDig, dig)
				args.SetBytesV("dig", string(hexDig))
			} else {
				// new account -- set their password
				args.Set("key", key)
			}
			statusCode, body, err := fasthttp.Post(nil, logIn, &args)
			if err != nil {
				return nil, err
			}
			switch statusCode {
			case 200, 201:
			case 404:
				msg = "It doesn't look like that was your password." // TODO: translate
				continue
			default:
				return nil, pkgErrors.Wrap(errors.New("HTTP "+strconv.Itoa(statusCode)), logIn)
			}
			token = new(Token)
			if err := json.Unmarshal(body, &token); err != nil {
				return nil, err
			}
			return token, nil
		}
		continue
	}
}

func register(host, email, username string) error { // TODO: autofill email and username
	// browser.OpenURL panics on unsupported operating systems
	defer func() {
		if err := recover(); err != nil {
			return err
		}
	}()
	if err := browser.OpenURL(path.Join(host, "register")); err != nil {
		return err
	}
	return nil
}
