package account

type TokenKey struct {
	Name string
}

func (t TokenKey) String() string {
	return t.Name
}

type TokenInfo struct {
	Token     string `json:"access_token"`
	Type      string `json:"token_type"`
	Expires   int    `json:"expires"`
	CreatedAt string `bson:"created_at" json:"created_at"`
}
