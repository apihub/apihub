package connection

type Connection interface{}

type connection struct {
	listenNetwork string
	listenAddr    string
}

func New(listenNetwork, listenAddr string) *connection {
	return &connection{
		listenNetwork: listenNetwork,
		listenAddr:    listenAddr,
	}
}
