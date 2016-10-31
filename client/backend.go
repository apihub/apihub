package client

import "github.com/apihub/apihub"

type backend struct {
	apihub.BackendInfo
}

func newBackend(info apihub.BackendInfo) *backend {
	return &backend{info}
}

func (b *backend) Info() apihub.BackendInfo {
	return b.BackendInfo
}

func (b *backend) Start() error {
	panic("not implemented")
}

func (b *backend) Stop() error {
	panic("not implemented")
}
