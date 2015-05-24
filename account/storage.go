package account

// Storage is an interface for "storage".
// To be compatible, the Storage which implements this interface must pass the acceptance suite that could be found
// in the folder storage/test/suite.go.
type Storage interface {
	//SaveToken inserts new content in the storage. i.e: User, TokenInfo.
	SaveToken(TokenKey, interface{}) error
	// GetToken returns the content in the storage by given key.
	GetToken(TokenKey) (interface{}, error)
	//DeleteToken deletes a content by given ley.
	DeleteToken(TokenKey) error
}
