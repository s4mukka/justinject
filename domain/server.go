package domain

import "github.com/gin-gonic/gin"

type IRouter interface {
	gin.IRouter
	Use(middleware ...gin.HandlerFunc) gin.IRoutes
	Run(addr ...string) (err error)
}

type IServer interface {
	Init(intializeRoutes func(router IRouter) error, port int) error
}

type IServerFactory interface {
	MakeServer(ctx IContext) IServer
}
