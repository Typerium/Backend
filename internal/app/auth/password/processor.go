package password

type Processor interface {
	Encode(password string) (string, error)
	Equal(hash, password string) (bool, error)
}
