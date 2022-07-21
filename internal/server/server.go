package server

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-playground/validator"
	"github.com/go-redis/redis/v8"
	"github.com/jmoiron/sqlx"
	"github.com/labstack/echo/v4"

	"github.com/dinorain/pinjembuku/config"
	"github.com/dinorain/pinjembuku/internal/middlewares"
	"github.com/dinorain/pinjembuku/pkg/logger"

	bookDeliveryHTTP "github.com/dinorain/pinjembuku/internal/book/delivery/http/handlers"
	librarianDeliveryHTTP "github.com/dinorain/pinjembuku/internal/librarian/delivery/http/handlers"
	orderDeliveryHTTP "github.com/dinorain/pinjembuku/internal/order/delivery/http/handlers"
	userDeliveryHTTP "github.com/dinorain/pinjembuku/internal/user/delivery/http/handlers"

	bookUseCase "github.com/dinorain/pinjembuku/internal/book/usecase"
	librarianUseCase "github.com/dinorain/pinjembuku/internal/librarian/usecase"
	orderUseCase "github.com/dinorain/pinjembuku/internal/order/usecase"
	sessUseCase "github.com/dinorain/pinjembuku/internal/session/usecase"
	userUseCase "github.com/dinorain/pinjembuku/internal/user/usecase"

	librarianRepository "github.com/dinorain/pinjembuku/internal/librarian/repository"
	orderRepository "github.com/dinorain/pinjembuku/internal/order/repository"
	sessRepository "github.com/dinorain/pinjembuku/internal/session/repository"
	userRepository "github.com/dinorain/pinjembuku/internal/user/repository"
)

type Server struct {
	logger      logger.Logger
	cfg         *config.Config
	v           *validator.Validate
	echo        *echo.Echo
	mw          middlewares.MiddlewareManager
	db          *sqlx.DB
	redisClient *redis.Client
}

// Server constructor
func NewAppServer(logger logger.Logger, cfg *config.Config, db *sqlx.DB, redisClient *redis.Client) *Server {
	return &Server{
		logger:      logger,
		cfg:         cfg,
		v:           validator.New(),
		echo:        echo.New(),
		db:          db,
		redisClient: redisClient,
	}
}

// Run service
func (s *Server) Run() error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	s.mw = middlewares.NewMiddlewareManager(s.logger, s.cfg)

	userRepo := userRepository.NewUserPGRepository(s.db)
	librarianRepo := librarianRepository.NewLibrarianPGRepository(s.db)
	orderRepo := orderRepository.NewOrderPGRepository(s.db)

	sessRepo := sessRepository.NewSessionRepository(s.redisClient, s.cfg)
	userRedisRepo := userRepository.NewUserRedisRepo(s.redisClient, s.logger)
	librarianRedisRepo := librarianRepository.NewLibrarianRedisRepo(s.redisClient, s.logger)
	orderRedisRepo := orderRepository.NewOrderRedisRepo(s.redisClient, s.logger)

	sessUC := sessUseCase.NewSessionUseCase(sessRepo, s.cfg)
	userUC := userUseCase.NewUserUseCase(s.cfg, s.logger, userRepo, userRedisRepo)
	librarianUC := librarianUseCase.NewLibrarianUseCase(s.cfg, s.logger, librarianRepo, librarianRedisRepo)
	bookUC := bookUseCase.NewBookUseCase(s.cfg, s.logger)
	orderUC := orderUseCase.NewOrderUseCase(s.cfg, s.logger, orderRepo, orderRedisRepo)

	l, err := net.Listen("tcp", s.cfg.Server.Port)
	if err != nil {
		return err
	}
	defer l.Close()

	userHandlers := userDeliveryHTTP.NewUserHandlersHTTP(s.echo.Group("user"), s.logger, s.cfg, s.mw, s.v, userUC, sessUC)
	userHandlers.UserMapRoutes()

	librarianHandlers := librarianDeliveryHTTP.NewLibrarianHandlersHTTP(s.echo.Group("librarian"), s.logger, s.cfg, s.mw, s.v, librarianUC, sessUC)
	librarianHandlers.LibrarianMapRoutes()

	bookHandlers := bookDeliveryHTTP.NewBookHandlersHTTP(s.echo.Group("book"), s.logger, s.cfg, s.mw, s.v, bookUC)
	bookHandlers.BookMapRoutes()

	orderHandlers := orderDeliveryHTTP.NewOrderHandlersHTTP(s.echo.Group("order"), s.logger, s.cfg, s.mw, s.v, orderUC, bookUC, userUC, librarianUC, sessUC)
	orderHandlers.OrderMapRoutes()

	go func() {
		if err := s.runHttpServer(); err != nil {
			s.logger.Errorf("s.runHttpServer: %v", err)
			cancel()
		}
	}()

	<-ctx.Done()
	if err := s.echo.Server.Shutdown(ctx); err != nil {
		s.logger.WarnMsg("echo.Server.Shutdown", err)
	}

	return nil
}
