package proto

import (
	validation "github.com/go-ozzo/ozzo-validation"
	"github.com/go-ozzo/ozzo-validation/is"
)

func (m *ProfilesUser) Validate() error {
	return validation.ValidateStruct(m,
		validation.Field(&m.Email, is.Email),
	)
}
