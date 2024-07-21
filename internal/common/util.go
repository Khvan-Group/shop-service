package common

import (
	"github.com/Khvan-Group/common-library/utils"
	"github.com/gorilla/context"
	"net/http"
)

func GetJwtUser(r *http.Request) JwtUser {
	return JwtUser{
		Login: utils.ToString(context.Get(r, "login")),
		Role:  utils.ToString(context.Get(r, "role")),
	}
}
