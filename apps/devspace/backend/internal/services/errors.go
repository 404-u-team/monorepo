package services

import "errors"

var ErrInternal = errors.New("внутренняя ошибка сервера")
var ErrUnauthorized = errors.New("нет доступа с данными реквизитами")
var ErrEmptyPayload = errors.New("все поля пустые, нечего изменять")

var ErrUserNotLeader = errors.New("пользователь не является лидером проекта")
var ErrUserExists = errors.New("пользователь с такой почтой уже зарегистрирован")
var ErrUserNotFound = errors.New("пользователь с такой почтой и/или паролем не найден")
var ErrUserConflict = errors.New("пользователь с таким nickname уже существует")

var ErrProjectConflict = errors.New("проект с таким названием уже существует")
var ErrProjectNotFound = errors.New("проект не найден")
var ErrProjectHasSlots = errors.New("проект не может быть удален из-за наличия связанных слотов")

var ErrSlotConflict = errors.New("слот с такими значениями уже существует")
var ErrSlotNotFound = errors.New("слот с таким ID не найден")
