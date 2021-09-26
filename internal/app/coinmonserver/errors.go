package coinmonserver

import "errors"

const (
	// STATUS_OK - Статус ok для json ответа
	STATUS_OK = "ok"

	// STATUS_ERROR - Статус error для json ответа
	STATUS_ERROR = "error"
)

var (
	//ErrInvalidAccessToken - Ошибка для невалидных токенов
	ErrInvalidAccessToken error

	//ErrTokenExpired - Ошибка для просроченных токенов
	ErrTokenExpired error

	//ErrUserDoesNotExist - Ошибка для несуществующих пользователей
	ErrUserDoesNotExist error
)

func init() {
	ErrInvalidAccessToken = errors.New("invalid auth token")
	ErrTokenExpired = errors.New("token is expired")
	ErrUserDoesNotExist = errors.New("user/pass not found")
}
