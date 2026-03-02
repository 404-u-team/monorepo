package db

import (
	"errors"
)

var UniqueKeyDuplErr error = errors.New("unique key(s) duplication error")
