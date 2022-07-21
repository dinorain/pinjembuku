package librarian

import "github.com/labstack/echo/v4"

// Librarian HTTP Handlers interface
type LibrarianHandlers interface {
	Register() echo.HandlerFunc
	Login() echo.HandlerFunc
	GetMe() echo.HandlerFunc
	FindAll() echo.HandlerFunc
	FindById() echo.HandlerFunc
	UpdateById() echo.HandlerFunc
	DeleteById() echo.HandlerFunc
	Logout() echo.HandlerFunc
	RefreshToken() echo.HandlerFunc
}
