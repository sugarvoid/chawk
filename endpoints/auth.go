package endpoints

import "fmt"

func GetToken() string {
	return fmt.Sprintf("/learn/api/public/v1/oauth2/token")
}
