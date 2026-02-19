package chawk

import "errors"

const MAX_RESPONSE_SIZE int64 = 3 * 1024 * 1024 // 3MB

var ErrEmptyStringParameter = errors.New("there is a blank parameter that is empty string")
