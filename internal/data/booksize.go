package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrInvalidBooksizeFormat = errors.New("invalid booksize format")

type Booksize int32

func (bs Booksize) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d pages", bs)

	quotedJsonValue := strconv.Quote(jsonValue)

	return []byte(quotedJsonValue), nil
}

func (bs *Booksize) UnmarshalJSON(jsonValue []byte) error {
	unquotedJsonValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidBooksizeFormat
	}

	parts := strings.Split(unquotedJsonValue, " ")
	if len(parts) != 2 || parts[1] != "pages" {
		return ErrInvalidBooksizeFormat
	}

	i, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidBooksizeFormat
	}

	*bs = Booksize(i)

	return nil
}
