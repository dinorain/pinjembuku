package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt"
	"github.com/golang/mock/gomock"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/stretchr/testify/require"

	"github.com/dinorain/pinjembuku/config"
	"github.com/dinorain/pinjembuku/internal/middlewares"
	"github.com/dinorain/pinjembuku/internal/models"
	"github.com/dinorain/pinjembuku/internal/librarian/delivery/http/dto"
	"github.com/dinorain/pinjembuku/internal/librarian/mock"
	mockSessUC "github.com/dinorain/pinjembuku/internal/session/mock"
	"github.com/dinorain/pinjembuku/pkg/converter"
	"github.com/dinorain/pinjembuku/pkg/logger"
)

func TestLibrariansHandler_Register(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	librarianUC := mock.NewMockLibrarianUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)

	appLogger := logger.NewAppLogger(nil)
	mw := middlewares.NewMiddlewareManager(appLogger, nil)

	e := echo.New()
	v := validator.New()
	cfg := &config.Config{Session: config.Session{Expire: 1234}}
	handlers := NewLibrarianHandlersHTTP(e.Group("librarian"), appLogger, cfg, mw, v, librarianUC, sessUC)

	reqDto := &dto.LibrarianRegisterRequestDto{
		Email:         "email@gmail.com",
		FirstName:     "FirstName",
		LastName:      "LastName",
		Password:      "123456",
	}

	buf := &bytes.Buffer{}
	_ = json.NewEncoder(buf).Encode(reqDto)

	req := httptest.NewRequest(http.MethodPost, "/librarian", buf)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	ctx := e.NewContext(req, res)

	resDto := &dto.LibrarianRegisterResponseDto{
		LibrarianID: uuid.Nil,
	}

	buf, _ = converter.AnyToBytesBuffer(resDto)

	librarianUC.EXPECT().Register(gomock.Any(), gomock.Any()).AnyTimes().Return(&models.Librarian{}, nil)
	require.NoError(t, handlers.Register()(ctx))
	require.Equal(t, http.StatusCreated, res.Code)
	require.Equal(t, buf.String(), res.Body.String())
}

func TestLibrariansHandler_Login(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	librarianUC := mock.NewMockLibrarianUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)

	appLogger := logger.NewAppLogger(nil)
	mw := middlewares.NewMiddlewareManager(appLogger, nil)

	e := echo.New()
	v := validator.New()
	cfg := &config.Config{Session: config.Session{Expire: 1234}}
	handlers := NewLibrarianHandlersHTTP(e.Group("librarian"), appLogger, cfg, mw, v, librarianUC, sessUC)

	reqDto := &dto.LibrarianLoginRequestDto{
		Email:    "email@gmail.com",
		Password: "123456",
	}

	var buf bytes.Buffer
	_ = json.NewEncoder(&buf).Encode(reqDto)

	req := httptest.NewRequest(http.MethodPost, "/librarian/login", &buf)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	ctx := e.NewContext(req, res)

	mockLibrarian := &models.Librarian{
		LibrarianID:      uuid.New(),
		Email:         "email@gmail.com",
		FirstName:     "FirstName",
		LastName:      "LastName",
		Password:      "123456",
	}

	librarianUC.EXPECT().Login(gomock.Any(), reqDto.Email, reqDto.Password).AnyTimes().Return(mockLibrarian, nil)
	sessUC.EXPECT().CreateSession(gomock.Any(), &models.Session{UserID: mockLibrarian.LibrarianID}, cfg.Session.Expire).AnyTimes().Return("s", nil)
	librarianUC.EXPECT().GenerateTokenPair(gomock.Any(), gomock.Any()).AnyTimes().Return("rt", "at", nil)
	require.NoError(t, handlers.Login()(ctx))
	require.Equal(t, http.StatusCreated, res.Code)
}

func TestLibrariansHandler_FindAll(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	librarianUC := mock.NewMockLibrarianUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)

	appLogger := logger.NewAppLogger(nil)
	mw := middlewares.NewMiddlewareManager(appLogger, nil)

	e := echo.New()
	v := validator.New()
	cfg := &config.Config{Session: config.Session{Expire: 1234}}
	handlers := NewLibrarianHandlersHTTP(e.Group("librarian"), appLogger, cfg, mw, v, librarianUC, sessUC)

	req := httptest.NewRequest(http.MethodGet, "/librarian", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	ctx := e.NewContext(req, res)

	var librarians []models.Librarian
	librarians = append(librarians, models.Librarian{
		LibrarianID:      uuid.New(),
		Email:         "email@gmail.com",
		FirstName:     "FirstName",
		LastName:      "LastName",
		Password:      "123456",
	})

	librarianUC.EXPECT().FindAll(gomock.Any(), gomock.Any()).AnyTimes().Return(librarians, nil)
	require.NoError(t, handlers.FindAll()(ctx))
	require.Equal(t, http.StatusOK, res.Code)
}

func TestLibrariansHandler_FindById(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	librarianUC := mock.NewMockLibrarianUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)

	cfg := &config.Config{Session: config.Session{Expire: 1234}}
	appLogger := logger.NewAppLogger(cfg)
	mw := middlewares.NewMiddlewareManager(appLogger, nil)

	e := echo.New()
	v := validator.New()
	handlers := NewLibrarianHandlersHTTP(e.Group("librarian"), appLogger, cfg, mw, v, librarianUC, sessUC)

	req := httptest.NewRequest(http.MethodGet, "/librarian/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	ctx := e.NewContext(req, res)

	ctx.SetParamNames("id")
	ctx.SetParamValues("2ceba62a-35f4-444b-a358-4b14834837e1")

	librarianUC.EXPECT().CachedFindById(gomock.Any(), gomock.Any()).AnyTimes().Return(&models.Librarian{}, nil)
	require.NoError(t, handlers.FindById()(ctx))
	require.Equal(t, http.StatusOK, res.Code)
}

func TestLibrariansHandler_UpdateById(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	librarianUC := mock.NewMockLibrarianUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)

	cfg := &config.Config{Session: config.Session{Expire: 1234}}
	appLogger := logger.NewAppLogger(cfg)
	mw := middlewares.NewMiddlewareManager(appLogger, cfg)

	e := echo.New()
	e.Use(middleware.JWT([]byte("secret")))
	v := validator.New()
	handlers := NewLibrarianHandlersHTTP(e.Group("librarian"), appLogger, cfg, mw, v, librarianUC, sessUC)

	change := "changed"
	reqDto := &dto.LibrarianUpdateRequestDto{
		FirstName:     &change,
		LastName:      &change,
		Password:      &change,
		Avatar:        &change,
	}

	buf := &bytes.Buffer{}
	_ = json.NewEncoder(buf).Encode(reqDto)

	librarianUUID := uuid.New()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["session_id"] = uuid.New().String()
	claims["librarian_id"] = librarianUUID.String()
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	validToken, _ := token.SignedString([]byte("secret"))

	req := httptest.NewRequest(http.MethodPost, "/librarian/:id", buf)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("bearer %v", validToken))

	t.Run("Forbidden update by other librarian", func(t *testing.T) {
		t.Parallel()

		res := httptest.NewRecorder()
		ctx := e.NewContext(req, res)

		handler := handlers.UpdateById()
		h := middleware.JWTWithConfig(middleware.JWTConfig{
			Claims:     claims,
			SigningKey: []byte("secret"),
		})(handler)

		ctx.SetParamNames("id")
		ctx.SetParamValues("2ceba62a-35f4-444b-a358-4b14834837e1")

		librarianUC.EXPECT().UpdateById(gomock.Any(), gomock.Any()).AnyTimes().Return(&models.Librarian{LibrarianID: librarianUUID}, nil)

		require.NoError(t, h(ctx))
		require.Equal(t, http.StatusForbidden, res.Code)
	})

	t.Run("Success", func(t *testing.T) {
		t.Parallel()

		res := httptest.NewRecorder()
		ctx := e.NewContext(req, res)

		handler := handlers.UpdateById()
		h := middleware.JWTWithConfig(middleware.JWTConfig{
			Claims:     claims,
			SigningKey: []byte("secret"),
		})(handler)

		ctx.SetParamNames("id")
		ctx.SetParamValues(librarianUUID.String())

		librarianUC.EXPECT().UpdateById(gomock.Any(), gomock.Any()).AnyTimes().Return(&models.Librarian{LibrarianID: librarianUUID}, nil)
		librarianUC.EXPECT().FindById(gomock.Any(), librarianUUID).AnyTimes().Return(&models.Librarian{LibrarianID: librarianUUID}, nil)

		require.NoError(t, h(ctx))
		require.Equal(t, http.StatusOK, res.Code)
	})
}

func TestLibrariansHandler_DeleteById(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	librarianUC := mock.NewMockLibrarianUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)

	cfg := &config.Config{Session: config.Session{Expire: 1234}}
	appLogger := logger.NewAppLogger(cfg)
	mw := middlewares.NewMiddlewareManager(appLogger, nil)

	e := echo.New()
	v := validator.New()
	handlers := NewLibrarianHandlersHTTP(e.Group("librarian"), appLogger, cfg, mw, v, librarianUC, sessUC)

	req := httptest.NewRequest(http.MethodDelete, "/librarian/:id", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	res := httptest.NewRecorder()
	ctx := e.NewContext(req, res)

	librarianUUID := uuid.New()
	ctx.SetParamNames("id")
	ctx.SetParamValues(librarianUUID.String())

	librarianUC.EXPECT().DeleteById(gomock.Any(), librarianUUID).AnyTimes().Return(nil)
	require.NoError(t, handlers.DeleteById()(ctx))
	require.Equal(t, http.StatusOK, res.Code)
}

func TestLibrariansHandler_GetMe(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	librarianUC := mock.NewMockLibrarianUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)

	cfg := &config.Config{Session: config.Session{Expire: 1234}}
	appLogger := logger.NewAppLogger(cfg)
	mw := middlewares.NewMiddlewareManager(appLogger, cfg)

	e := echo.New()
	e.Use(middleware.JWT([]byte("secret")))
	v := validator.New()
	handlers := NewLibrarianHandlersHTTP(e.Group("librarian"), appLogger, cfg, mw, v, librarianUC, sessUC)

	librarianUUID := uuid.New()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["session_id"] = uuid.New().String()
	claims["librarian_id"] = librarianUUID.String()
	claims["role"] = "librarian"
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	validToken, _ := token.SignedString([]byte("secret"))

	req := httptest.NewRequest(http.MethodPost, "/librarian/logout", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("bearer %v", validToken))

	res := httptest.NewRecorder()
	ctx := e.NewContext(req, res)

	handler := handlers.GetMe()
	h := middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     claims,
		SigningKey: []byte("secret"),
	})(handler)

	sessUC.EXPECT().GetSessionById(gomock.Any(), claims["session_id"].(string)).AnyTimes().Return(&models.Session{}, nil)
	librarianUC.EXPECT().CachedFindById(gomock.Any(), gomock.Any()).AnyTimes().Return(&models.Librarian{}, nil)

	require.NoError(t, h(ctx))
	require.Equal(t, http.StatusOK, res.Code)
}

func TestLibrariansHandler_Logout(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	librarianUC := mock.NewMockLibrarianUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)

	cfg := &config.Config{Session: config.Session{Expire: 1234}}
	appLogger := logger.NewAppLogger(cfg)
	mw := middlewares.NewMiddlewareManager(appLogger, cfg)

	e := echo.New()
	e.Use(middleware.JWT([]byte("secret")))
	v := validator.New()
	handlers := NewLibrarianHandlersHTTP(e.Group("librarian"), appLogger, cfg, mw, v, librarianUC, sessUC)

	librarianUUID := uuid.New()
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["session_id"] = uuid.New().String()
	claims["librarian_id"] = librarianUUID.String()
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()
	validToken, _ := token.SignedString([]byte("secret"))

	req := httptest.NewRequest(http.MethodPost, "/librarian/logout", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	req.Header.Set(echo.HeaderAuthorization, fmt.Sprintf("bearer %v", validToken))

	res := httptest.NewRecorder()
	ctx := e.NewContext(req, res)

	handler := handlers.Logout()
	h := middleware.JWTWithConfig(middleware.JWTConfig{
		Claims:     claims,
		SigningKey: []byte("secret"),
	})(handler)

	sessUC.EXPECT().DeleteById(gomock.Any(), claims["session_id"].(string)).AnyTimes().Return(nil)

	require.NoError(t, h(ctx))
	require.Equal(t, http.StatusOK, res.Code)
}

func TestLibrariansHandler_RefreshToken(t *testing.T) {
	t.Parallel()

	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	librarianUC := mock.NewMockLibrarianUseCase(ctrl)
	sessUC := mockSessUC.NewMockSessUseCase(ctrl)

	cfg := &config.Config{Session: config.Session{Expire: 1234}, Server: config.ServerConfig{JwtSecretKey: "secret"}}
	appLogger := logger.NewAppLogger(cfg)
	mw := middlewares.NewMiddlewareManager(appLogger, cfg)

	e := echo.New()
	v := validator.New()
	handlers := NewLibrarianHandlersHTTP(e.Group("librarian"), appLogger, cfg, mw, v, librarianUC, sessUC)

	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["session_id"] = uuid.New().String()
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()
	validToken, _ := token.SignedString([]byte("secret"))

	reqDto := &dto.LibrarianRefreshTokenDto{
		RefreshToken: validToken,
	}

	buf := &bytes.Buffer{}
	_ = json.NewEncoder(buf).Encode(reqDto)

	req := httptest.NewRequest(http.MethodPost, "/librarian/refresh", buf)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

	res := httptest.NewRecorder()
	ctx := e.NewContext(req, res)

	sessUC.EXPECT().GetSessionById(gomock.Any(), claims["session_id"].(string)).AnyTimes().Return(&models.Session{}, nil)
	librarianUC.EXPECT().FindById(gomock.Any(), gomock.Any()).AnyTimes().Return(&models.Librarian{}, nil)
	librarianUC.EXPECT().GenerateTokenPair(gomock.Any(), gomock.Any()).AnyTimes().Return("rt", "at", nil)

	require.NoError(t, handlers.RefreshToken()(ctx))
	require.Equal(t, http.StatusOK, res.Code)
}
