package api

import (
	"net/http"

	"github.com/Brian3647/nakme/pkg/api/auth"
	"github.com/Brian3647/nakme/pkg/db"
	"github.com/Brian3647/nakme/pkg/util"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func AddRoutes(e *echo.Echo) {
	// Middleware
	e.Use(middleware.CORS())
	//e.Use(middleware.HTTPSRedirect())

	// Rate limiter
	if util.InProd() {
		rl := middleware.NewRateLimiterMemoryStore(10)
		e.Use(middleware.RateLimiter(rl))
	}

	err := db.SetUp(); if err != nil {
		panic(err)
	}

	util.SetUp()

	// Routes
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, World!\n")
	})

	e.GET("/api", func(c echo.Context) error {
		return c.JSON(http.StatusOK, c.Echo().Routes())
	})

	e.GET("/health", Health)
	e.GET("/api/health", Health)

	e.POST("/api/auth/signup", auth.SignUp)
	e.POST("/api/auth/confirm_token", auth.ConfirmToken)
	e.GET("/api/auth/confirm_email", auth.ConfirmEmail)
	e.POST("/api/auth/login", auth.LogIn)
	e.POST("/api/auth/delete", auth.DeleteUser)
}
