package common

type JwtUser struct {
	Login string
	Role  string
}

type Wallet struct {
	User  string `json:"username" database:"username"`
	Total int    `json:"total" database:"total"`
}

type WalletUpdate struct {
	Total    int    `json:"total" database:"total"`
	Username string `json:"username" database:"username"`
	Action   string `json:"action"`
}

const (
	WALLET_TOTAL_ADD       = "ADD"
	WALLET_TOTAL_SUBSTRUCT = "SUBSTRUCT"
)
