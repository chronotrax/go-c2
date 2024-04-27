package handler

import (
	"fmt"
)

func newParseError(err error) map[string]string {
	return map[string]string{"error": fmt.Sprintf("failed to parse request: %s", err)}
}

var internalServerError = map[string]string{"error": "internal server error"}
