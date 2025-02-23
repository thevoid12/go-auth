package routes

import (
	"context"
	"goauth/web/handler"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func Initialize(ctx context.Context, l *zap.Logger) (router *gin.Engine) {
	l.Sugar().Info("Initializing logger")

	router = gin.Default()
	router.Use(gin.Recovery())

	//auth group sets the context and calls auth middleware
	//	rAuth := router.Group("/auth")
	// rAuth.Use(middleware.ContextMiddleware(ctx), middleware.AuthMiddleware(ctx))
	rWeb := router.Group("/web")
	rWeb.GET("/login", handler.LoginPageHandler)
	rWeb.GET("/auth/google", handler.GoogleauthHandler)
	rWeb.GET("/callback", handler.CallbackHandler)
	rWeb.GET("/home", handler.HomePageHandler)
	for _, route := range router.Routes() {
		l.Sugar().Infof("Route: %s %s", route.Method, route.Path)
	}

	return router
}
