package system

import "fmt"

type EnvError struct {
	VarName string
}

type OSError struct {
	Err string
}

type ParseError struct{}

func (error_object EnvError) Error() string {
	return fmt.Sprintf("Ошибка получения переченной окружения %s", error_object.VarName)
}

func (error_object OSError) Error() string {
	return fmt.Sprintf("Ошибка ОС: %s", error_object.Err)
}

func (error_object ParseError) Error() string {
	return fmt.Sprintf("Ошибка парсинга конфига")
}
