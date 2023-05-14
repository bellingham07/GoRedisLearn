package myRouters

import (
	"GoRedisLearn/controller"
	"GoRedisLearn/middleware"
	"github.com/gin-gonic/gin"
)

type UserRoute struct {
}

func (*UserRoute) InitRouter(g *gin.RouterGroup) {
	u := g.Group("/user")
	u.Use(middleware.RefreshTokenMiddleware())
	// 不需要拦截器
	{
		u.POST("/code", controller.Code)
		u.POST("/login", controller.Login)
		u.POST("/test", controller.Test)
	}
	// 需要拦截器
	c := g.Group("/user")
	c.Use(middleware.RefreshTokenMiddleware(), middleware.AuthMiddleware())
	{
		c.POST("/id", controller.QueryById)
	}
}