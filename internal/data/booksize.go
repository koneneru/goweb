package data

import (
	"fmt"
	"strconv"
)

type Booksize int32

func (bs Booksize) MarshalJSON() ([]byte, error) {
	jsonValue := fmt.Sprintf("%d pages", bs)

	quotedJsonValue := strconv.Quote(jsonValue)

	return []byte(quotedJsonValue), nil
}
