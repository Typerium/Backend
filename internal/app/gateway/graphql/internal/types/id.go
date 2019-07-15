package types

import (
	"io"

	"github.com/99designs/gqlgen/graphql"
	"github.com/pkg/errors"
)

// ID graphql model identifier
type ID struct {
	input    interface{}
	str      string
	isString bool
	num      int
	isNumber bool
}

// StringToID create graphql identifier from string identifier
func StringToID(in string) *ID {
	return &ID{
		isString: true,
		str:      in,
	}
}

// Int32ToID create graphql identifier from int32 identifier
func Int32ToID(in int32) *ID {
	return &ID{
		isNumber: true,
		num:      int(in),
	}
}

// UnmarshalGQL decode from graphql query to model
func (id *ID) UnmarshalGQL(v interface{}) error {
	id.input = v
	return nil
}

// MarshalGQL encode from model to graphql response
func (id ID) MarshalGQL(w io.Writer) {
	if id.isString {
		graphql.MarshalID(id.str).MarshalGQL(w)
		return
	}

	if id.isNumber {
		graphql.MarshalIntID(id.num).MarshalGQL(w)
		return
	}
}

// String convert graphql identifier to string
func (id *ID) String() (out string, err error) {
	if id.isString {
		return id.str, nil
	}

	id.str, err = graphql.UnmarshalID(id.input)
	if err != nil {
		return "", errors.WithStack(err)
	}

	id.isString = true
	out = id.str
	return
}

// Int32 convert graphql identifier to int32 type
func (id *ID) Int32() (out int32, err error) {
	if id.isNumber {
		return int32(id.num), nil
	}

	id.num, err = graphql.UnmarshalIntID(id.input)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	id.isNumber = true
	out = int32(id.num)
	return
}
