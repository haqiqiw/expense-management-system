package route

import (
	"embed"
	internalHttp "expense-management-system/internal/delivery/http"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
)

//go:embed embeds/openapi.json
var swaggerJSON []byte

//go:embed embeds/swagger-ui/*
var swaggerUI embed.FS

type RouteConfig struct {
	App            *gin.Engine
	AuthMiddlware  gin.HandlerFunc
	AuthController *internalHttp.AuthController
	UserController *internalHttp.UserController
}

func (c *RouteConfig) Setup() {
	SetupSwagger(c.App)

	c.App.GET("/", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "OK")
	})
	c.App.GET("/health", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "OK")
	})

	api := c.App.Group("/api")

	// without auth
	api.POST("/login", c.AuthController.Login)
	api.POST("/users", c.UserController.Register)

	// with auth
	api.POST("/logout", c.AuthMiddlware, c.AuthController.Logout)
	api.GET("/users/me", c.AuthMiddlware, c.UserController.Me)
}

func SetupSwagger(app *gin.Engine) {
	subFS, err := fs.Sub(swaggerUI, "embeds/swagger-ui")
	if err != nil {
		panic("failed to create swagger-ui sub-filesystem: " + err.Error())
	}

	app.StaticFS("/swagger", http.FS(subFS))

	app.GET("/api/openapi.json", func(ctx *gin.Context) {
		ctx.Header("Content-Type", "application/json")
		ctx.Data(http.StatusOK, "application/json", swaggerJSON)
	})
}
