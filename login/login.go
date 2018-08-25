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

type Login struct {
	var Email, Username string
	var PasswordHash    []byte
}

type Account struct {
	ID          int64  `json:"id"`
	RowID       int64  `json:"rowid"`
	HasPassword bool   `json:"has_pw"`
}

func LogIn(host string) error {
	var login Login
	var msg, password string
	var button int8
	user := path.Join(host, "user")
	for {
		dialog.New("email", &login.Email, &username, &button, msg)
		switch button {
		case dialog.Cancel:
			return error911.NewCancel()
		case dialog.Register:
			if err := register(host, login.Email, login.Username); err != nil {
				return err
			}
			continue
		}
		var args fasthttp.Args
		statusCode, body, err := fasthttp.Post(nil, user, &args)
		if err != nil {
			return err
		}
		switch statusCode {
		case 200:
		case 404:
			msg = "E-mail or username not found." // TODO: translate
			continue
		default:
			return pkgErrors.Wrap(errors.New("HTTP "+strconv.Itoa(statusCode)), user)
		}

		var account Account
		if err := json.Unmarshal(body, &account); err != nil {
			return err
		}
		if account.HasPassword {
			// old account -- ask for their password
		} else {
			// new account -- create a password
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
