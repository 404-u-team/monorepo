package db

type UniqueKeyDuplErr struct{}

func (error_object UniqueKeyDuplErr) Error() string {
	return "unique key(s) duplication error"
}
