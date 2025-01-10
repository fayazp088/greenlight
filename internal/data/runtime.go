package data

import (
	"fmt"
	"strconv"
)

type Runtime int32

func (r Runtime) MarshalJSON() ([]byte, error) {

	jsonValue := fmt.Sprintf("%d min", r)

	quotedJson := strconv.Quote(jsonValue)

	return []byte(quotedJson), nil
}
