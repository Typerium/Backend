package web

import (
	"sync"

	"github.com/valyala/fasthttp"
)

type ClientFactory interface {
	Acquire() *fasthttp.Client
	Release(client *fasthttp.Client)
}

type fasthttpClientFactory struct {
	pool *sync.Pool
}

func NewFasthttpClientFactory() ClientFactory {
	return &fasthttpClientFactory{
		pool: &sync.Pool{
			New: func() interface{} {
				return &fasthttp.Client{
				}
			},
		},
	}
}

func (f *fasthttpClientFactory) Acquire() *fasthttp.Client {
	return f.pool.Get().(*fasthttp.Client)
}

func (f *fasthttpClientFactory) Release(client *fasthttp.Client) {
	f.pool.Put(client)
}
