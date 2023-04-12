package api

type CreateAccountResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type LoginResponse struct {
	JWT     string `json:"jwt-token"`
	Refresh string `json:"refresh-token"`
}

type RefreshTokenResponse struct {
	JWT     string `json:"jwt-token"`
	Refresh string `json:"refresh-token"`
}

type GetAccountResponse struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type GetAllAccountsResponse []*GetAccountResponse

type DeleteAccountResponse struct {
	Deleted int `json:"deleted"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
