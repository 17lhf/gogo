package router

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"

	"gogo/internal/cache"
	"gogo/internal/config"
	"gogo/internal/handler"
	"gogo/internal/i18n"
	"gogo/internal/middleware"
	"gogo/internal/pkg"
	"gogo/internal/repository"
)

// Dependencies holds all application components for route setup.
type Dependencies struct {
	AuthHandler     *handler.AuthHandler
	UserHandler     *handler.UserHandler
	RoleHandler     *handler.RoleHandler
	MenuHandler     *handler.MenuHandler
	StoreHandler    *handler.StoreHandler
	TerminalHandler *handler.TerminalHandler
	LogHandler      *handler.LogHandler
	StatsHandler    *handler.StatsHandler

	SessionCache   *cache.SessionCache
	HeartbeatCache *cache.HeartbeatCache
	UserRepo       repository.UserRepository
	LogRepo        repository.LogRepository
	Enforcer       *casbin.Enforcer
	AuthConfig     config.AuthConfig
}

// Register sets up all routes and middleware.
func Register(r *gin.Engine, d *Dependencies) {
	r.Use(middleware.RequestLog())
	r.Use(i18n.Middleware())

	// Health check (no auth)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	api := r.Group("/api/v1")

	// Auth routes (public and protected)
	auth := api.Group("/auth")
	{
		auth.POST("/login", d.AuthHandler.Login)

		// Protected auth routes
		authUse := auth.Group("")
		authUse.Use(middleware.Auth(d.SessionCache, d.AuthConfig))
		{
			authUse.POST("/logout", d.AuthHandler.Logout)
			authUse.GET("/me", d.AuthHandler.Me, middleware.Audit(d.LogRepo))
			authUse.PUT("/password", d.AuthHandler.ChangePassword)
		}
	}

	// All protected routes below
	protected := api.Group("")
	protected.Use(middleware.Auth(d.SessionCache, d.AuthConfig))
	protected.Use(middleware.PasswordExpiry(d.UserRepo, d.AuthConfig))
	protected.Use(middleware.Permission(d.Enforcer))
	protected.Use(middleware.DataScope(d.UserRepo))

	// User routes (with audit)
	protectedUsers := protected.Group("")
	protectedUsers.Use(middleware.Audit(d.LogRepo))
	{
		protectedUsers.GET("/users", d.UserHandler.List)
		protectedUsers.GET("/users/:id", d.UserHandler.GetByID)
		protectedUsers.POST("/users", d.UserHandler.Create)
		protectedUsers.PUT("/users/:id", d.UserHandler.Update)
		protectedUsers.DELETE("/users/:id", d.UserHandler.Delete)
		protectedUsers.PUT("/users/:id/password", d.UserHandler.ResetPassword)
		protectedUsers.PUT("/users/:id/roles", d.UserHandler.AssignRoles)
		protectedUsers.PUT("/users/:id/stores", d.UserHandler.AssignStores)
	}

	// Role routes
	protectedRoles := protected.Group("")
	protectedRoles.Use(middleware.Audit(d.LogRepo))
	{
		protectedRoles.GET("/roles", d.RoleHandler.List)
		protectedRoles.GET("/roles/:id", d.RoleHandler.GetByID)
		protectedRoles.POST("/roles", d.RoleHandler.Create)
		protectedRoles.PUT("/roles/:id", d.RoleHandler.Update)
		protectedRoles.DELETE("/roles/:id", d.RoleHandler.Delete)
		protectedRoles.GET("/roles/:id/menus", d.RoleHandler.GetMenus)
		protectedRoles.PUT("/roles/:id/menus", d.RoleHandler.AssignMenus)
	}

	// Menu routes
	protectedMenus := protected.Group("")
	protectedMenus.Use(middleware.Audit(d.LogRepo))
	{
		protectedMenus.GET("/menus", d.MenuHandler.Tree)
		protectedMenus.GET("/menus/:id", d.MenuHandler.GetByID)
		protectedMenus.POST("/menus", d.MenuHandler.Create)
		protectedMenus.PUT("/menus/:id", d.MenuHandler.Update)
		protectedMenus.DELETE("/menus/:id", d.MenuHandler.Delete)
	}

	// Store routes
	protectedStores := protected.Group("")
	protectedStores.Use(middleware.Audit(d.LogRepo))
	{
		protectedStores.GET("/stores", d.StoreHandler.List)
		protectedStores.GET("/stores/:id", d.StoreHandler.GetByID)
		protectedStores.POST("/stores", d.StoreHandler.Create)
		protectedStores.PUT("/stores/:id", d.StoreHandler.Update)
		protectedStores.DELETE("/stores/:id", d.StoreHandler.Delete)
	}

	// Terminal routes
	protectedTerminals := protected.Group("")
	protectedTerminals.Use(middleware.Audit(d.LogRepo))
	{
		protectedTerminals.GET("/terminals", d.TerminalHandler.List)
		protectedTerminals.GET("/terminals/:id", d.TerminalHandler.GetByID)
		protectedTerminals.POST("/terminals", d.TerminalHandler.Create)
		protectedTerminals.PUT("/terminals/:id", d.TerminalHandler.Update)
		protectedTerminals.DELETE("/terminals/:id", d.TerminalHandler.Delete)
	}

	// Terminal device routes (authenticated via X-Device-Token, no JWT)
	terminal := api.Group("/terminals")
	{
		terminal.POST("/:sn/heartbeat", d.TerminalHandler.Heartbeat)
		terminal.POST("/:sn/rotate-token", d.TerminalHandler.RotateToken)
	}

	// Log routes
	protectedLogs := protected.Group("")
	{
		protectedLogs.GET("/logs/operations", d.LogHandler.ListOperations)
		protectedLogs.GET("/logs/terminals", d.LogHandler.ListTerminals)
	}

	// Stats routes
	protectedStats := protected.Group("")
	{
		protectedStats.GET("/stats/terminals", d.StatsHandler.GetTerminals)
		protectedStats.GET("/stats/users", d.StatsHandler.GetUsers)
	}

	// 404 handler
	r.NoRoute(func(c *gin.Context) {
		pkg.Error(c, 404, pkg.CodeParamError, i18n.Localize(c, i18n.MsgEndpointNotFound))
	})
}
