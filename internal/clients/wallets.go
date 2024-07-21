package clients

import (
	"encoding/json"
	"github.com/Khvan-Group/common-library/constants"
	"github.com/Khvan-Group/common-library/errors"
	"github.com/Khvan-Group/common-library/utils"
	"github.com/go-resty/resty/v2"
	"net/http"
	"shop-service/internal/common"
)

var WALLET_SERVICE_URL string

func GetWalletByUser(username string, client *resty.Client) (*common.Wallet, *errors.CustomError) {
	var result common.Wallet
	WALLET_SERVICE_URL = utils.GetEnv("WALLET_SERVICE_URL")
	request := client.R()
	request.Header.Set(constants.X_IS_INTERNAL_SERVICE, "true")

	response, err := request.Get(WALLET_SERVICE_URL + "/wallets/" + username)
	if err != nil {
		return nil, errors.NewInternal("Внутренняя ошибка: Возможно сервис кошельков не доступен.")
	}

	if response.StatusCode() != http.StatusOK {
		return nil, errors.NewBadRequest("Ошибка получения кошелька пользователя.")
	}

	if err = json.Unmarshal(response.Body(), &result); err != nil {
		return nil, errors.NewInternal("Failed unmarshalling wallet response body")
	}

	return &result, nil
}
