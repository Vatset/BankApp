package handler

import (
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	_ "users-balance-monitoring/docs"
	"users-balance-monitoring/pkg/service"
)

type Handler struct {
	service *service.Service
}

func NewHandler(service *service.Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) InitRoutes() *gin.Engine {

	router := gin.New()

	ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("http://localhost:8080/swagger/doc.json"),
		ginSwagger.DefaultModelsExpandDepth(-1))

	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp)
		auth.POST("/sign-in", h.signIn)
	}
	api_user := router.Group("/api", h.userIdentity)
	{
		balance := api_user.Group("/balance")
		{
			balance.GET("/", h.getBalance)
			balance.GET("/history", h.getAllTransactions)
			balance.POST("/transfer", h.balanceTransfer)
		}
	}
	api_admin := router.Group("/api_admin")
	{
		balance := api_admin.Group("/balance")
		{
			balance.POST("/deposit", h.balanceDeposit)
			balance.POST("/withdraw", h.balanceWithdraw)
		}
	}
	return router
}
