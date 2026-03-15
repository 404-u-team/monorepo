package services

import "errors"

var ErrInternal = errors.New("внутренняя ошибка сервера")
var ErrUnauthorized = errors.New("нет доступа с данными реквизитами")
var ErrEmptyPayload = errors.New("все поля пустые, нечего изменять")
var ErrBadRequest = errors.New("")

var ErrUserNotLeader = errors.New("пользователь не является лидером проекта")
var ErrUserLeader = errors.New("пользователь не может выполнить это действие, так как является лидером проекта")
var ErrUserExists = errors.New("пользователь с такой почтой уже зарегистрирован")
var ErrUserNotFound = errors.New("пользователь не найден")
var ErrUserConflict = errors.New("пользователь с таким nickname уже существует")

var ErrProjectConflict = errors.New("проект с таким названием уже существует")
var ErrProjectNotFound = errors.New("проект не найден")
var ErrProjectHasSlots = errors.New("проект не может быть удален из-за наличия связанных слотов")

var ErrSlotConflict = errors.New("слот с такими значениями уже существует")
var ErrSlotNotFound = errors.New("слот с таким ID не найден")
var ErrSlotIsClosed = errors.New("нельзя создать запрос для закрытого слота")

var ErrRowAlreadyExists = errors.New("такая запись уже существует")
var ErrRowNotExists = errors.New("такой записи не существует")

var ErrProjectRequestConflict = errors.New("такая заявка на проект уже существует")
var ErrCantInviteYourself = errors.New("нельзя пригласить самого себя")
