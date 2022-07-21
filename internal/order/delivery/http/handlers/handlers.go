package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/go-playground/validator"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/dinorain/pinjembuku/config"
	"github.com/dinorain/pinjembuku/internal/book"
	"github.com/dinorain/pinjembuku/internal/librarian"
	"github.com/dinorain/pinjembuku/internal/middlewares"
	"github.com/dinorain/pinjembuku/internal/models"
	"github.com/dinorain/pinjembuku/internal/order"
	"github.com/dinorain/pinjembuku/internal/order/delivery/http/dto"
	"github.com/dinorain/pinjembuku/internal/session"
	"github.com/dinorain/pinjembuku/internal/user"
	"github.com/dinorain/pinjembuku/pkg/constants"
	httpErrors "github.com/dinorain/pinjembuku/pkg/http_errors"
	"github.com/dinorain/pinjembuku/pkg/logger"
	"github.com/dinorain/pinjembuku/pkg/utils"
)

type orderHandlersHTTP struct {
	group       *echo.Group
	logger      logger.Logger
	cfg         *config.Config
	mw          middlewares.MiddlewareManager
	v           *validator.Validate
	orderUC     order.OrderUseCase
	bookUC      book.BookUseCase
	userUC      user.UserUseCase
	librarianUC librarian.LibrarianUseCase
	sessUC      session.SessUseCase
}

var _ order.OrderHandlers = (*orderHandlersHTTP)(nil)

func NewOrderHandlersHTTP(
	group *echo.Group,
	logger logger.Logger,
	cfg *config.Config,
	mw middlewares.MiddlewareManager,
	v *validator.Validate,
	orderUC order.OrderUseCase,
	bookUC book.BookUseCase,
	userUC user.UserUseCase,
	librarianUC librarian.LibrarianUseCase,
	sessUC session.SessUseCase,
) *orderHandlersHTTP {
	return &orderHandlersHTTP{group: group, logger: logger, cfg: cfg, mw: mw, v: v, orderUC: orderUC, userUC: userUC, librarianUC: librarianUC, bookUC: bookUC, sessUC: sessUC}
}

// Create
// @Tags Orders
// @Summary To register order
// @Description Admin create order
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param payload body dto.OrderCreateRequestDto true "Payload"
// @Success 200 {object} dto.OrderCreateResponseDto
// @Router /order [post]
func (h *orderHandlersHTTP) Create() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		createDto := &dto.OrderCreateRequestDto{}
		if err := c.Bind(createDto); err != nil {
			h.logger.WarnMsg("bind", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		if err := h.v.StructCtx(ctx, createDto); err != nil {
			h.logger.WarnMsg("validate", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		sessID, _, _, err := h.getSessionIDFromCtx(c)
		if err != nil {
			h.logger.Errorf("getSessionIDFromCtx: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		session, err := h.sessUC.GetSessionById(ctx, sessID)
		if err != nil {
			h.logger.Errorf("sessUC.GetSessionById: %v", err)
			if errors.Is(err, redis.Nil) {
				return httpErrors.NewUnauthorizedError(c, nil, h.cfg.Http.DebugErrorsResponse)
			}
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		user, err := h.userUC.CachedFindById(ctx, session.UserID)
		if err != nil {
			h.logger.Errorf("userUC.CachedFindById: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		book, err := h.bookUC.FindByWork(ctx, createDto.BookKey)
		if err != nil {
			h.logger.Errorf("bookUC.FindByWork: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		order, err := h.registerReqToOrderModel(createDto, user, nil, book)
		if err != nil {
			h.logger.Errorf("orderHandlersHTTP.registerReqToOrderModel: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		createdOrder, err := h.orderUC.Create(ctx, order)
		if err != nil {
			h.logger.Errorf("orderUC.Create: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		return c.JSON(http.StatusCreated, dto.OrderCreateResponseDto{OrderID: createdOrder.OrderID})
	}
}

// FindAll
// @Tags Orders
// @Summary Find all orders
// @Description Find all orders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param size query string false "pagination size"
// @Param page query string false "pagination page"
// @Success 200 {object} dto.OrderFindResponseDto
// @Router /order [get]
func (h *orderHandlersHTTP) FindAll() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		pq := utils.NewPaginationFromQueryParams(c.QueryParam(constants.Size), c.QueryParam(constants.Page))

		var orders []models.Order
		_, userID, role, err := h.getSessionIDFromCtx(c)
		if err != nil {
			h.logger.Errorf("getSessionIDFromCtx: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}
		if role == models.UserRoleUser {
			userUUID, _ := uuid.Parse(userID)
			if res, err := h.orderUC.FindAllByUserId(ctx, userUUID, pq); err != nil {
				h.logger.Errorf("orderUC.FindAllByUserId: %v", err)
				return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
			} else {
				orders = res
			}
		} else {
			if res, err := h.orderUC.FindAll(ctx, pq); err != nil {
				h.logger.Errorf("orderUC.FindAll: %v", err)
				return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
			} else {
				orders = res
			}
		}

		return c.JSON(http.StatusOK, dto.OrderFindResponseDto{
			Data: orders,
			Meta: utils.PaginationMetaDto{
				Limit:  pq.GetLimit(),
				Offset: pq.GetOffset(),
				Page:   pq.GetPage(),
			},
		})
	}
}

// FindById
// @Tags Orders
// @Summary Find order
// @Description Find existing order by id
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} dto.OrderResponseDto
// @Router /order/{id} [get]
func (h *orderHandlersHTTP) FindById() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		orderUUID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			h.logger.WarnMsg("uuid.FromString", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		order, err := h.orderUC.CachedFindById(ctx, orderUUID)
		if err != nil {
			h.logger.Errorf("orderUC.CachedFindById: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		return c.JSON(http.StatusOK, dto.OrderResponseFromModel(order))
	}
}

// AcceptById
// @Tags Orders
// @Summary Accept order
// @Description Librarian accept order
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Order ID"
// @Success 200 {object} dto.OrderResponseDto
// @Router /order/{id} [put]
func (h *orderHandlersHTTP) AcceptById() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		orderUUID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			h.logger.WarnMsg("uuid.FromString", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		updateDto := &dto.OrderUpdateRequestDto{
			Status: models.OrderStatusAccepted,
		}

		order, err := h.orderUC.FindById(ctx, orderUUID)
		if err != nil {
			h.logger.Errorf("orderUC.FindById: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		_, userID, role, err := h.getSessionIDFromCtx(c)
		if err != nil {
			h.logger.Errorf("getSessionIDFromCtx: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		if role != "" {
			return httpErrors.NewForbiddenError(c, nil, h.cfg.Http.DebugErrorsResponse)
		}

		order, err = h.updateReqToOrderModel(order, updateDto)
		if err != nil {
			h.logger.Errorf("orderHandlersHTTP.updateReqToOrderModel: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		if librarianUUID, err := uuid.Parse(userID); err != nil {
			h.logger.Errorf("uuid.Parse: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		} else {
			order.LibrarianID = &librarianUUID
		}
		order, err = h.orderUC.UpdateById(ctx, order)
		if err != nil {
			h.logger.Errorf("orderUC.UpdateById: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		return c.JSON(http.StatusOK, dto.OrderResponseFromModel(order))
	}
}

func (h *orderHandlersHTTP) getSessionIDFromCtx(c echo.Context) (sessionID string, userID string, role string, err error) {
	user, ok := c.Get("user").(*jwt.Token)
	if !ok {
		h.logger.Warnf("jwt.Token: %+v", c.Get("user"))
		return "", "", "", errors.New("invalid token header")
	}

	claims, ok := user.Claims.(jwt.MapClaims)
	if !ok {
		h.logger.Warnf("jwt.MapClaims: %+v", c.Get("user"))
		return "", "", "", errors.New("invalid token header")
	}

	sessionID, ok = claims["session_id"].(string)
	if !ok {
		h.logger.Warnf("session_id: %+v", claims)
		return "", "", "", errors.New("invalid token header")
	}

	role, _ = claims["role"].(string)
	if role != "" {
		userID, _ = claims["user_id"].(string)
	} else {
		userID, _ = claims["librarian_id"].(string)
	}
	return sessionID, userID, role, nil
}

func (h *orderHandlersHTTP) registerReqToOrderModel(r *dto.OrderCreateRequestDto, user *models.User, librarian *models.Librarian, book *models.Book) (*models.Order, error) {
	var librarianID *uuid.UUID
	if librarian != nil {
		librarianID = &librarian.LibrarianID
	}

	orderCandidate := &models.Order{
		UserID:      user.UserID,
		LibrarianID: librarianID,
		Item: models.OrderItem{
			BookKey:         book.BookKey,
			Title:           book.Title,
			EditionCount:    book.EditionCount,
			CoverID:         book.CoverID,
			CoverEditionKey: book.CoverEditionKey,
		},
		Status:         models.OrderStatusPending,
		PickupSchedule: r.PickupSchedule,
	}

	return orderCandidate, nil
}

func (h *orderHandlersHTTP) updateReqToOrderModel(updateCandidate *models.Order, r *dto.OrderUpdateRequestDto) (*models.Order, error) {
	if r.Status != models.OrderStatusPending && r.Status != models.OrderStatusAccepted {
		return nil, fmt.Errorf("status invalid: %v", r.Status)
	}

	if r.Status == updateCandidate.Status {
		return nil, fmt.Errorf("order has been %v", r.Status)
	}
	updateCandidate.Status = r.Status

	return updateCandidate, nil
}
