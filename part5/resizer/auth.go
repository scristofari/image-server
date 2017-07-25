package resizer

import "fmt"

var credentialsList = map[string]string{
	"app1": "passApp1",
}

func CheckCredentials(user string, password string) error {
	pass, ok := credentialsList[user]
	if !ok {
		return fmt.Errorf("user invalid")
	}
	if pass != password {
		return fmt.Errorf("password invalid")
	}
	return nil
}
