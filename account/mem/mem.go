// package mem provides in memory storage implementation, for test purposes.
package mem

import (
	. "github.com/backstage/backstage/account"
	"github.com/backstage/backstage/errors"
)

type Mem struct {
	Tokens map[TokenKey]interface{}
}

func New() Storable {
	return &Mem{
		Tokens: map[TokenKey]interface{}{},
	}
}

func (m *Mem) SaveToken(k TokenKey, d interface{}) error {
	m.Tokens[k] = d
	return nil
}

func (m *Mem) GetToken(k TokenKey) (interface{}, error) {
	t, ok := m.Tokens[k]
	if !ok {
		return nil, &errors.NotFoundError{Payload: "Token not found."}
	}
	return t, nil
}

func (m *Mem) DeleteToken(k TokenKey) error {
	if _, ok := m.Tokens[k]; !ok {
		return &errors.NotFoundError{Payload: "Token not found."}
	}
	delete(m.Tokens, k)
	return nil
}
