package db

import (
	"errors"
)

var ErrUniqueKeyDupl error = errors.New("unique key(s) duplication error")
