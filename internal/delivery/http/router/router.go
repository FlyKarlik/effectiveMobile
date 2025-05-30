package http_router

import (
	"net/http/pprof"

	http_handler "github.com/FlyKarlik/effectiveMobile/internal/delivery/http/handler"
	http_middleware "github.com/FlyKarlik/effectiveMobile/internal/delivery/http/middleware"

	_ "github.com/FlyKarlik/effectiveMobile/docs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type HTTPRouter struct {
	middleware *http_middleware.HTTPMiddleware
	handler    *http_handler.HTTPHandler
}

func New(middleware *http_middleware.HTTPMiddleware, handler *http_handler.HTTPHandler) *HTTPRouter {
	return &HTTPRouter{
		middleware: middleware,
		handler:    handler,
	}
}

func (h *HTTPRouter) InitRouter() *gin.Engine {
	router := gin.New()
	router.Use(gin.Logger())
	router.Use(gin.Recovery())
	router.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowCredentials: true,
		AllowMethods:     []string{"*"},
		AllowHeaders:     []string{"*"},
	}))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.GET("/ping", h.handler.Ping)
	registerPprof(router)

	api := router.Group("v1", h.middleware.JSONMiddleware())
	{
		h.registerUserRoutes(api)
	}

	return router
}

func (h *HTTPRouter) registerUserRoutes(router *gin.RouterGroup) {
	userGroup := router.Group("users")
	{
		userGroup.GET("/", h.handler.SearchUsers)
		userGroup.POST("/", h.handler.CreateUser)
		userGroup.PATCH("/:id", h.handler.UpdateUser)
		userGroup.DELETE("/:id", h.handler.DeleteUser)
	}
}

func registerPprof(router *gin.Engine) {
	pprofGroup := router.Group("/debug/pprof")
	{
		router.GET("/debug/vars", gin.WrapH(pprof.Handler("vars")))
		pprofGroup.GET("/", gin.WrapF(pprof.Index))
		pprofGroup.GET("/cmdline", gin.WrapF(pprof.Cmdline))
		pprofGroup.GET("/profile", gin.WrapF(pprof.Profile))
		pprofGroup.POST("/symbol", gin.WrapF(pprof.Symbol))
		pprofGroup.GET("/symbol", gin.WrapF(pprof.Symbol))
		pprofGroup.GET("/trace", gin.WrapF(pprof.Trace))
		pprofGroup.GET("/allocs", gin.WrapH(pprof.Handler("allocs")))
		pprofGroup.GET("/block", gin.WrapH(pprof.Handler("block")))
		pprofGroup.GET("/goroutine", gin.WrapH(pprof.Handler("goroutine")))
		pprofGroup.GET("/heap", gin.WrapH(pprof.Handler("heap")))
		pprofGroup.GET("/mutex", gin.WrapH(pprof.Handler("mutex")))
		pprofGroup.GET("/threadcreate", gin.WrapH(pprof.Handler("threadcreate")))
	}
}
