package signature

type KeyCreator interface {
	Acquire() (Key, error)
	Release(key Key)
	Create(key []byte) (Key, error)
}

type Key interface {
	Key() interface{}
	Bytes() []byte
}
