// package mem provides in memory storage implementation, for test purposes.
package mem

import (
	"encoding/json"

	"github.com/backstage/backstage/account"
)

type Mem struct {
	Tokens map[account.TokenKey][]byte
}

func New() account.Storable {
	return &Mem{
		Tokens: map[account.TokenKey][]byte{},
	}
}

func (m *Mem) SaveToken(k account.TokenKey, expires int, d interface{}) error {
	data, _ := json.Marshal(d)
	m.Tokens[k] = data
	return nil
}

func (m *Mem) GetToken(k account.TokenKey, t interface{}) error {
	if token, ok := m.Tokens[k]; ok {
		if err := json.Unmarshal(token, &t); err != nil {
			panic(err)
		}
	}
	return nil
}

func (m *Mem) DeleteToken(k account.TokenKey) error {
	delete(m.Tokens, k)
	return nil
}
