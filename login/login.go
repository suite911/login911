package login

import (
	"github.com/valyala/fasthttp"
)

func LogIn() {
	var email, username, password string
	var button int8
	for {
		dialog.New("email", &email, &username, &button)
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
