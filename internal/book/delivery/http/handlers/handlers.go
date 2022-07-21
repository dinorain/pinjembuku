package handlers

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/labstack/echo/v4"

	"github.com/dinorain/pinjembuku/config"
	"github.com/dinorain/pinjembuku/internal/book"
	"github.com/dinorain/pinjembuku/internal/book/delivery/http/dto"
	"github.com/dinorain/pinjembuku/internal/middlewares"
	"github.com/dinorain/pinjembuku/pkg/constants"
	httpErrors "github.com/dinorain/pinjembuku/pkg/http_errors"
	"github.com/dinorain/pinjembuku/pkg/logger"
	"github.com/dinorain/pinjembuku/pkg/utils"
)

type bookHandlersHTTP struct {
	group  *echo.Group
	logger logger.Logger
	cfg    *config.Config
	mw     middlewares.MiddlewareManager
	v      *validator.Validate
	bookUC book.BookUseCase
}

var _ book.BookHandlers = (*bookHandlersHTTP)(nil)

func NewBookHandlersHTTP(
	group *echo.Group,
	logger logger.Logger,
	cfg *config.Config,
	mw middlewares.MiddlewareManager,
	v *validator.Validate,
	bookUC book.BookUseCase,
) *bookHandlersHTTP {
	return &bookHandlersHTTP{group: group, logger: logger, cfg: cfg, mw: mw, v: v, bookUC: bookUC}
}

// FindBySubject
// @Tags Books
// @Summary Find all books of certain subject
// @Description Find all books of certain subject
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param size query string false "pagination size"
// @Param page query string false "pagination page"
// @Success 200 {object} dto.BookFindResponseDto
// @Router /book [get]
func (h *bookHandlersHTTP) FindBySubject() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		subject := c.Param("subject")
		if subject == "" {
			return httpErrors.ErrorCtxResponse(c, errors.New("subject required"), h.cfg.Http.DebugErrorsResponse)
		}

		pq := utils.NewPaginationFromQueryParams(c.QueryParam(constants.Size), c.QueryParam(constants.Page))

		books, err := h.bookUC.FindAllBySubject(ctx, subject, pq)
		if err != nil {
			h.logger.Errorf("userUC.FindAll: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		var data []dto.BookReponseDto
		for _, book := range books {
			data = append(data, dto.BookReponseDto{
				BookKey:         book.BookKey,
				Title:           book.Title,
				EditionCount:    book.EditionCount,
				CoverID:         book.CoverID,
				CoverEditionKey: book.CoverEditionKey,
			})
		}

		return c.JSON(http.StatusOK, dto.BookFindResponseDto{
			Data: data,
			Meta: utils.PaginationMetaDto{
				Limit:  pq.GetLimit(),
				Offset: pq.GetOffset(),
				Page:   pq.GetPage(),
			},
		})
	}
}
