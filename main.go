package main

import (
	"context"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/gin-gonic/gin"

	"gogo/internal/cache"
	"gogo/internal/casbin"
	"gogo/internal/config"
	"gogo/internal/db"
	"gogo/internal/handler"
	"gogo/internal/model"
	"gogo/internal/repository"
	"gogo/internal/router"
	"gogo/internal/service"
)

func main() {
	var level slog.Level
	switch os.Getenv("LOG_LEVEL") {
	case "debug":
		level = slog.LevelDebug
	case "warn":
		level = slog.LevelWarn
	case "error":
		level = slog.LevelError
	default:
		level = slog.LevelInfo
	}

	var h slog.Handler
	if os.Getenv("LOG_FORMAT") == "text" {
		h = slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	} else {
		h = slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: level})
	}
	slog.SetDefault(slog.New(h))

	cfg := config.Load()

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// PostgreSQL pool (pgx)
	pool, err := db.New(ctx, cfg.Postgres.DSN())
	if err != nil {
		log.Fatalf("postgres init: %v", err)
	}
	defer pool.Close()

	if err := db.Ping(ctx, pool); err != nil {
		log.Fatalf("postgres unreachable: %v", err)
	}
	slog.Info("postgres connected")

	// GORM
	gormDB, err := db.NewGORM(ctx, cfg.Postgres.DSN())
	if err != nil {
		log.Fatalf("gorm init: %v", err)
	}

	// Redis
	rdb := cache.New(cfg.Redis.Addr, cfg.Redis.Password)
	defer rdb.Close()

	if err := cache.Ping(context.Background(), rdb); err != nil {
		log.Fatalf("redis unreachable: %v", err)
	}
	slog.Info("redis connected")

	// Casbin
	enforcer, err := casbin.NewEnforcer(gormDB)
	if err != nil {
		log.Fatalf("casbin init: %v", err)
	}
	slog.Info("casbin initialized")

	// Repositories
	userRepo := repository.NewUserRepository(gormDB)
	roleRepo := repository.NewRoleRepository(gormDB)
	menuRepo := repository.NewMenuRepository(gormDB)
	storeRepo := repository.NewStoreRepository(gormDB)
	terminalRepo := repository.NewTerminalRepository(gormDB)
	logRepo := repository.NewLogRepository(gormDB)

	// Cache services
	sessionCache := cache.NewSessionCache(rdb, cfg.Auth.SessionTTL)
	lockoutCache := cache.NewLockoutCache(rdb, cfg.Auth.LockoutThreshold, cfg.Auth.LockoutDuration)
	heartbeatCache := cache.NewHeartbeatCache(rdb, 60*time.Second)

	// Services
	authSvc := service.NewAuthService(userRepo, roleRepo, sessionCache, lockoutCache, cfg.Auth)
	userSvc := service.NewUserService(userRepo, sessionCache)
	roleSvc := service.NewRoleService(roleRepo, menuRepo, enforcer)
	menuSvc := service.NewMenuService(menuRepo)
	storeSvc := service.NewStoreService(storeRepo)
	terminalSvc := service.NewTerminalService(terminalRepo, storeRepo, heartbeatCache, logRepo)

	// Handlers
	authHandler := handler.NewAuthHandler(authSvc)
	userHandler := handler.NewUserHandler(userSvc)
	roleHandler := handler.NewRoleHandler(roleSvc)
	menuHandler := handler.NewMenuHandler(menuSvc)
	storeHandler := handler.NewStoreHandler(storeSvc)
	terminalHandler := handler.NewTerminalHandler(terminalSvc, heartbeatCache)
	logHandler := handler.NewLogHandler(logRepo)

	// Start heartbeat expiry listener
	go cache.ListenForExpiry(context.Background(), rdb, func(sn string) {
		slog.Info("heartbeat expired", "sn", sn)
		terminalSvc.HandleStatusTimeout(context.Background(), sn)
	})

	// Register custom validators
	model.RegisterValidators()

	// Gin router
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())

	router.Register(r, &router.Dependencies{
		AuthHandler:     authHandler,
		UserHandler:     userHandler,
		RoleHandler:     roleHandler,
		MenuHandler:     menuHandler,
		StoreHandler:    storeHandler,
		TerminalHandler: terminalHandler,
		LogHandler:      logHandler,

		SessionCache:   sessionCache,
		HeartbeatCache: heartbeatCache,
		UserRepo:       userRepo,
		LogRepo:        logRepo,
		Enforcer:       enforcer,
		AuthConfig:     cfg.Auth,
	})

	slog.Info("HTTP server starting", "port", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		log.Fatalf("server: %v", err)
	}
}
