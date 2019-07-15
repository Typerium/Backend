package types

import (
	"fmt"
	"io"
	"strconv"
	"strings"
	"time"

	"github.com/pkg/errors"
)

// NewTimezone constructor for Timezone
func NewTimezone(name string, offset int32) *Timezone {
	return &Timezone{
		Name:     name,
		Offset:   offset,
		Location: time.FixedZone(name, int(offset)),
	}
}

// DefaultTimezone default of timezone type
var DefaultTimezone = &Timezone{
	Name:     "UTC",
	Location: time.UTC,
}

// Timezone type for implement graphql scalar type TimeZone
type Timezone struct {
	Name     string
	Offset   int32
	Location *time.Location
}

// UnmarshalGQL decode from graphql query to model
func (t *Timezone) UnmarshalGQL(input interface{}) (err error) {
	switch input := input.(type) {
	case string:
		posSpace := strings.Index(input, " ")
		if posSpace <= 0 {
			return errors.New("not found name of timezone")
		}
		t.Name = input[:posSpace]
		offsetStr := input[posSpace+1:]
		if len(offsetStr) == 0 {
			return errors.New("not found offset of timezone")
		}
		sign := offsetStr[0]
		switch sign {
		case '+':
			t.Offset, err = parseOffset(offsetStr[1:])
		case '-':
			t.Offset, err = parseOffset(offsetStr[1:])
			t.Offset = -t.Offset
		default:
			t.Offset, err = parseOffset(offsetStr)
		}
		if err != nil {
			return
		}
		t.Location = time.FixedZone(t.Name, int(t.Offset))
		return
	default:
		return ErrWrongType
	}
}

func parseOffset(input string) (offset int32, err error) {
	var hours, minutes, seconds int
	switch len(input) {
	case 0:
		return
	case 1:
		hours, err = strconv.Atoi(input)
		input = ""
	default:
		hours, err = strconv.Atoi(input[:2])
		input = input[2:]
	}
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	switch len(input) {
	case 0:
		break
	case 1:
		minutes, err = strconv.Atoi(input)
	default:
		minutes, err = strconv.Atoi(input[:2])
		input = input[2:]
	}
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	switch len(input) {
	case 0:
		break
	case 1:
		seconds, err = strconv.Atoi(input)
	default:
		seconds, err = strconv.Atoi(input[:2])
		input = input[2:]
	}
	if err != nil {
		err = errors.WithStack(err)
		return
	}

	offset = int32(hours*360 + minutes*60 + seconds)
	return
}

// MarshalGQL encode from model to graphql response
func (t Timezone) MarshalGQL(w io.Writer) {
	sign := "+"
	if t.Offset < 0 {
		sign = "-"
		t.Offset = -t.Offset
	}
	fmt.Fprintf(w, "%s %s%d", t.Name, sign, t.Offset)
}
