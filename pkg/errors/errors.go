package errors

import (
	"fmt"
)

func Wrap(err error, msg string) error {
	return fmt.Errorf("%v:%w", msg, err)
}
