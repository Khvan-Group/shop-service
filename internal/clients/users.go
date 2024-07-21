package clients

import (
	"encoding/json"
	"github.com/Khvan-Group/common-library/constants"
	"github.com/Khvan-Group/common-library/utils"
	"github.com/go-resty/resty/v2"
)

var AUTH_SERVICE_URL string

func ExistsUser(username string, client *resty.Client) bool {
	AUTH_SERVICE_URL = utils.GetEnv("AUTH_SERVICE_URL")
	var exists bool
	request := client.R()
	request.Header.Set(constants.X_IS_INTERNAL_SERVICE, "true")

	response, err := request.Get(AUTH_SERVICE_URL + "/common/" + username)
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(response.Body(), &exists); err != nil {
		panic(err)
	}

	return exists
}
