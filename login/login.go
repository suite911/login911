package login

import (
	"github.com/suite911/error911"

	"github.com/pkg/browser"
	"github.com/valyala/fasthttp"
)

type Login struct {
	var Email, Username string
	var PasswordHash    []byte
}

func LogIn(host string) error {
	var login Login
	var password string
	var button int8
	for {
		dialog.New("email", &login.Email, &username, &button)
		switch button {
		case dialog.Cancel:
			return error911.NewCancel()
		case dialog.Register:
			if err := register(login.Email, login.Username); err != nil {
				return err
			}
			continue
		}
		var args fasthttp.Args
		statusCode, body, err := fasthttp.Post(nil, url, &args)
		if err != nil {
			return err
		}
		// look up

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
