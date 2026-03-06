package services

import "errors"

var ErrInternal = errors.New("внутренняя ошибка сервера")
var ErrUserExists = errors.New("пользователь с такой почтой уже зарегистрирован")
var ErrUserNotFound = errors.New("пользователь с такой почтой и/или паролем не найден")
var ErrUnauthorized = errors.New("нет доступа с данными реквизитами")

var ErrProjectConflict = errors.New("проект с таким названием уже существует")
var ErrProjectNotFound = errors.New("проект не найден")
var ErrProjectHasSlots = errors.New("проект не может быть удален из-за наличия связанных слотов")
