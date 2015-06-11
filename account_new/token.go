package account_new

type TokenInfo struct {
	Token     string `json:"access_token"`
	CreatedAt string `json:"created_at"`
	Expires   int    `json:"expires"`
	Type      string `json:"token_type"`
	User      *User  `json:"-"`
}
