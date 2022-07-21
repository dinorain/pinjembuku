package book

import "github.com/labstack/echo/v4"

// Book HTTP Handlers interface
type BookHandlers interface {
	FindBySubject() echo.HandlerFunc
}
