package account_new

type TokenInfo struct {
	CreatedAt string `json:"created_at"`
	Token     string `json:"access_token"`
	Expires   int    `json:"expires"`
	Type      string `json:"token_type"`
	User      *User  `json:"-"`
}
