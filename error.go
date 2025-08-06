package runtime

import "errors"

var ErrServiceNotExist = errors.New("service not found")
var ErrServiceMultiple = errors.New("service multiple")
