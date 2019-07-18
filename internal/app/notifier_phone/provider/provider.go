package provider

type Provider interface {
	Send(number string, body string) error
}
