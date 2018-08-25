package login

import (
	"encoding/json"
	"errors"
	"path"
	"strconv"

	"github.com/suite911/error911"

	"github.com/pkg/browser"
	pkgErrors "github.com/pkg/errors"
	"github.com/valyala/fasthttp"
)

type Account struct {
	ID          int64  `json:"id"`
	RowID       int64  `json:"rowid"`
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
	user := path.Join(host, "user")
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
				return nil, error911.NewCancel()
			}
			if account.HasPassword {
				// old account -- verify their password
			} else {
				// new account -- set their password
				if retype != password {
					continue
				}
			}
			var args fasthttp.Args
			statusCode, body, err := fasthttp.Post(nil, user, &args)
			if err != nil {
				return nil, err
			}
			switch statusCode {
			case 200, 201:
			case 404:
				msg = "E-mail or username not found." // TODO: translate
				continue
			default:
				return nil, pkgErrors.Wrap(errors.New("HTTP "+strconv.Itoa(statusCode)), user)
			}
		}

		password = ""
		button = dialog.Cancel
		dialog.New("log in", &email, &password, &username, &button)
		switch button {
		case dialog.Cancel:
			return
		case dialog.LogIn:
			fmt.Println("Logging in as (\""+email+"\", \""+username+"\")")
		case dialog.Register:
			if err := browser.OpenURL("https://localhost:10443/register"); err != nil {
				panic(err)
			}
		}
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
