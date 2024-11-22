package CADDY_FILE_SERVER

import (
	"errors"
	"fmt"
)

type AbortRequestError struct {
	Msg string
}

func (e *AbortRequestError) Error() string {
	return fmt.Sprintf("request aborted: %s", e.Msg)
}

var BypassRequestError = errors.New("bypass request")
