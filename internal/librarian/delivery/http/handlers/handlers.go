package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator"
	"github.com/go-redis/redis/v8"
	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/dinorain/pinjembuku/config"
	"github.com/dinorain/pinjembuku/internal/librarian"
	"github.com/dinorain/pinjembuku/internal/librarian/delivery/http/dto"
	"github.com/dinorain/pinjembuku/internal/middlewares"
	"github.com/dinorain/pinjembuku/internal/models"
	"github.com/dinorain/pinjembuku/internal/session"
	"github.com/dinorain/pinjembuku/pkg/constants"
	httpErrors "github.com/dinorain/pinjembuku/pkg/http_errors"
	"github.com/dinorain/pinjembuku/pkg/logger"
	"github.com/dinorain/pinjembuku/pkg/utils"
)

type librarianHandlersHTTP struct {
	group    *echo.Group
	logger   logger.Logger
	cfg      *config.Config
	mw       middlewares.MiddlewareManager
	v        *validator.Validate
	librarianUC librarian.LibrarianUseCase
	sessUC   session.SessUseCase
}

var _ librarian.LibrarianHandlers = (*librarianHandlersHTTP)(nil)

func NewLibrarianHandlersHTTP(
	group *echo.Group,
	logger logger.Logger,
	cfg *config.Config,
	mw middlewares.MiddlewareManager,
	v *validator.Validate,
	librarianUC librarian.LibrarianUseCase,
	sessUC session.SessUseCase,
) *librarianHandlersHTTP {
	return &librarianHandlersHTTP{group: group, logger: logger, cfg: cfg, mw: mw, v: v, librarianUC: librarianUC, sessUC: sessUC}
}

// Register
// @Tags Librarians
// @Summary To register librarian
// @Description Admin create librarian
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param payload body dto.LibrarianRegisterRequestDto true "Payload"
// @Success 200 {object} dto.LibrarianRegisterResponseDto
// @Router /librarian [post]
func (h *librarianHandlersHTTP) Register() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		createDto := &dto.LibrarianRegisterRequestDto{}
		if err := c.Bind(createDto); err != nil {
			h.logger.WarnMsg("bind", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		if err := h.v.StructCtx(ctx, createDto); err != nil {
			h.logger.WarnMsg("validate", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		librarian, err := h.registerReqToLibrarianModel(createDto)
		if err != nil {
			h.logger.Errorf("registerReqToLibrarianModel: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		createdLibrarian, err := h.librarianUC.Register(ctx, librarian)
		if err != nil {
			h.logger.Errorf("librarianUC.Register: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		return c.JSON(http.StatusCreated, dto.LibrarianRegisterResponseDto{LibrarianID: createdLibrarian.LibrarianID})
	}
}

// Login
// @Tags Librarians
// @Summary Librarian login
// @Description Librarian login with email and password
// @Accept json
// @Produce json
// @Param payload body dto.LibrarianLoginRequestDto true "Payload"
// @Success 200 {object} dto.LibrarianLoginResponseDto
// @Router /librarian/login [post]
func (h *librarianHandlersHTTP) Login() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		loginDto := &dto.LibrarianLoginRequestDto{}
		if err := c.Bind(loginDto); err != nil {
			h.logger.WarnMsg("bind", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		if err := h.v.StructCtx(ctx, loginDto); err != nil {
			h.logger.WarnMsg("validate", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		email := loginDto.Email
		if !utils.ValidateEmail(email) {
			h.logger.Errorf("ValidateEmail: %v", email)
			return httpErrors.ErrorCtxResponse(c, errors.New("invalid email"), h.cfg.Http.DebugErrorsResponse)
		}

		librarian, err := h.librarianUC.Login(ctx, email, loginDto.Password)
		if err != nil {
			h.logger.Errorf("librarianUC.Login: %v", email)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		session, err := h.sessUC.CreateSession(ctx, &models.Session{
			UserID: librarian.LibrarianID,
		}, h.cfg.Session.Expire)
		if err != nil {
			h.logger.Errorf("sessUC.CreateSession: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		accessToken, refreshToken, err := h.librarianUC.GenerateTokenPair(librarian, session)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusCreated, dto.LibrarianLoginResponseDto{LibrarianID: librarian.LibrarianID, Tokens: &dto.LibrarianRefreshTokenResponseDto{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		}})
	}
}

// FindAll
// @Tags Librarians
// @Summary Find all librarians
// @Description Admin find all librarians
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param size query string false "pagination size"
// @Param page query string false "pagination page"
// @Success 200 {object} dto.LibrarianFindResponseDto
// @Router /librarian [get]
func (h *librarianHandlersHTTP) FindAll() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		pq := utils.NewPaginationFromQueryParams(c.QueryParam(constants.Size), c.QueryParam(constants.Page))
		librarians, err := h.librarianUC.FindAll(ctx, pq)
		if err != nil {
			h.logger.Errorf("librarianUC.FindAll: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		return c.JSON(http.StatusOK, dto.LibrarianFindResponseDto{
			Data: librarians,
			Meta: utils.PaginationMetaDto{
				Limit:  pq.GetLimit(),
				Offset: pq.GetOffset(),
				Page:   pq.GetPage(),
			},
		})
	}
}

// FindById
// @Tags Librarians
// @Summary Find librarian
// @Description Find existing librarian by id
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} dto.LibrarianResponseDto
// @Router /librarian/{id} [get]
func (h *librarianHandlersHTTP) FindById() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		librarianUUID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			h.logger.WarnMsg("uuid.FromString", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		librarian, err := h.librarianUC.CachedFindById(ctx, librarianUUID)
		if err != nil {
			h.logger.Errorf("librarianUC.CachedFindById: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		return c.JSON(http.StatusOK, dto.LibrarianResponseFromModel(librarian))
	}
}

// UpdateById
// @Tags Librarians
// @Summary Update librarian
// @Description Update existing librarian
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "Librarian ID"
// @Param payload body dto.LibrarianUpdateRequestDto true "Payload"
// @Success 200 {object} dto.LibrarianResponseDto
// @Router /librarian/{id} [put]
func (h *librarianHandlersHTTP) UpdateById() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		librarianUUID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			h.logger.WarnMsg("uuid.FromString", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		_, librarianID, role, err := h.getSessionIDFromCtx(c)
		if err != nil {
			h.logger.Errorf("getSessionIDFromCtx: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		if role != models.UserRoleAdmin && librarianID != librarianUUID.String() {
			return httpErrors.NewForbiddenError(c, nil, h.cfg.Http.DebugErrorsResponse)
		}

		updateDto := &dto.LibrarianUpdateRequestDto{}
		if err := c.Bind(updateDto); err != nil {
			h.logger.WarnMsg("bind", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		if err := h.v.StructCtx(ctx, updateDto); err != nil {
			h.logger.WarnMsg("validate", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		librarian, err := h.librarianUC.FindById(ctx, librarianUUID)
		if err != nil {
			h.logger.Errorf("librarianUC.FindById: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		librarian, err = h.updateReqToLibrarianModel(librarian, updateDto)
		if err != nil {
			h.logger.Errorf("updateReqToLibrarianModel: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		librarian, err = h.librarianUC.UpdateById(ctx, librarian)
		if err != nil {
			h.logger.Errorf("librarianUC.UpdateById: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		return c.JSON(http.StatusOK, dto.LibrarianResponseFromModel(librarian))
	}
}

// DeleteById
// @Tags Librarians
// @Summary Delete librarian
// @Description Delete existing librarian, admin only
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} nil
// @Param id path string true "Librarian ID"
// @Router /librarian/{id} [delete]
func (h *librarianHandlersHTTP) DeleteById() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		librarianUUID, err := uuid.Parse(c.Param("id"))
		if err != nil {
			h.logger.WarnMsg("uuid.FromString", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		if err := h.librarianUC.DeleteById(ctx, librarianUUID); err != nil {
			h.logger.Errorf("librarianUC.DeleteById: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		return c.JSON(http.StatusOK, nil)
	}
}

// GetMe
// @Tags Librarians
// @Summary Find me
// @Description Get session id from token, find librarian by uuid and returns it
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} dto.LibrarianResponseDto
// @Router /librarian/me [get]
func (h *librarianHandlersHTTP) GetMe() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
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

		librarian, err := h.librarianUC.CachedFindById(ctx, session.UserID)
		if err != nil {
			h.logger.Errorf("librarianUC.CachedFindById: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		return c.JSON(http.StatusOK, dto.LibrarianResponseFromModel(librarian))
	}
}

// Logout
// @Tags Librarians
// @Summary Librarian logout
// @Description Delete current session
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} nil
// @Router /librarian/logout [post]
func (h *librarianHandlersHTTP) Logout() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()
		sessID, _, _, err := h.getSessionIDFromCtx(c)
		if err != nil {
			h.logger.Errorf("getSessionIDFromCtx: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		if err := h.sessUC.DeleteById(ctx, sessID); err != nil {
			h.logger.Errorf("sessUC.DeleteById: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		return c.JSON(http.StatusOK, nil)
	}
}

// RefreshToken
// @Tags Librarians
// @Summary Refresh access token
// @Description Refresh access token
// @Accept json
// @Produce json
// @Param payload body dto.LibrarianRefreshTokenDto true "Payload"
// @Success 200 {object} dto.LibrarianRefreshTokenResponseDto
// @Router /librarian/refresh [post]
func (h *librarianHandlersHTTP) RefreshToken() echo.HandlerFunc {
	return func(c echo.Context) error {
		ctx := c.Request().Context()

		refreshTokenDto := &dto.LibrarianRefreshTokenDto{}
		if err := c.Bind(refreshTokenDto); err != nil {
			h.logger.WarnMsg("bind", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		token, err := jwt.Parse(refreshTokenDto.RefreshToken, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				h.logger.Errorf("jwt.SigningMethodHMAC: %v", token.Header["alg"])
				return nil, fmt.Errorf("jwt.SigningMethodHMAC: %v", token.Header["alg"])
			}

			return []byte(h.cfg.Server.JwtSecretKey), nil
		})

		if err != nil {
			h.logger.Warnf("jwt.Parse")
			return httpErrors.ErrorCtxResponse(c, errors.New("invalid refresh token"), h.cfg.Http.DebugErrorsResponse)
		}

		if !token.Valid {
			h.logger.Warnf("token.Valid")
			return httpErrors.ErrorCtxResponse(c, errors.New("invalid refresh token"), h.cfg.Http.DebugErrorsResponse)
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			h.logger.Warnf("jwt.MapClaims: %+v", token.Claims)
			return httpErrors.ErrorCtxResponse(c, errors.New("invalid refresh token"), h.cfg.Http.DebugErrorsResponse)
		}

		sessID, ok := claims["session_id"].(string)
		if !ok {
			h.logger.Warnf("session_id: %+v", claims)
			return httpErrors.ErrorCtxResponse(c, errors.New("invalid refresh token"), h.cfg.Http.DebugErrorsResponse)
		}

		session, err := h.sessUC.GetSessionById(ctx, sessID)
		if err != nil {
			h.logger.Errorf("sessUC.GetSessionById: %v", err)
			if errors.Is(err, redis.Nil) {
				return httpErrors.NewUnauthorizedError(c, nil, h.cfg.Http.DebugErrorsResponse)
			}
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		librarian, err := h.librarianUC.FindById(ctx, session.UserID)
		if err != nil {
			h.logger.Errorf("librarianUC.FindById: %v", err)
			return httpErrors.ErrorCtxResponse(c, err, h.cfg.Http.DebugErrorsResponse)
		}

		accessToken, refreshToken, err := h.librarianUC.GenerateTokenPair(librarian, sessID)
		if err != nil {
			return err
		}

		return c.JSON(http.StatusOK, dto.LibrarianRefreshTokenResponseDto{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
		})
	}
}

func (h *librarianHandlersHTTP) getSessionIDFromCtx(c echo.Context) (sessionID string, librarianID string, role string, err error) {
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

	librarianID, _ = claims["librarian_id"].(string)
	role, _ = claims["role"].(string)

	return sessionID, librarianID, role, nil
}

func (h *librarianHandlersHTTP) registerReqToLibrarianModel(r *dto.LibrarianRegisterRequestDto) (*models.Librarian, error) {
	librarianCandidate := &models.Librarian{
		Email:         r.Email,
		FirstName:     r.FirstName,
		LastName:      r.LastName,
		Avatar:        nil,
		Password:      r.Password,
	}

	if err := librarianCandidate.PrepareCreate(); err != nil {
		return nil, err
	}

	return librarianCandidate, nil
}

func (h *librarianHandlersHTTP) updateReqToLibrarianModel(updateCandidate *models.Librarian, r *dto.LibrarianUpdateRequestDto) (*models.Librarian, error) {

	if r.FirstName != nil {
		updateCandidate.FirstName = strings.TrimSpace(*r.FirstName)
	}
	if r.LastName != nil {
		updateCandidate.LastName = strings.TrimSpace(*r.LastName)
	}
	if r.Avatar != nil {
		avatar := strings.TrimSpace(*r.Avatar)
		updateCandidate.Avatar = &avatar
	}
	if r.Password != nil {
		updateCandidate.Password = *r.Password
		if err := updateCandidate.HashPassword(); err != nil {
			return nil, err
		}
	}

	return updateCandidate, nil
}
