package router

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/wxlbd/polaris/internal/infrastructure/config"
	"github.com/wxlbd/polaris/internal/interface/http/handler"
	"github.com/wxlbd/polaris/internal/interface/middleware"
)

// NewRouter 创建并配置路由 (去家庭化架构)
// NewRouter 创建并配置路由
func NewRouter(
	cfg *config.Config,
	authHandler *handler.AuthHandler,
	uploadHandler *handler.UploadHandler,
	logger *zap.Logger,
) *gin.Engine {
	// 设置Gin运行模式
	gin.SetMode(cfg.Server.Mode)

	r := gin.New()

	// 全局中间件
	r.Use(middleware.CORS())
	r.Use(middleware.Logger())
	r.Use(gin.Recovery())

	// 静态文件服务
	r.Static("/uploads", "./uploads")

	// API v1 路由组
	v1 := r.Group("/v1")
	{
		// 认证相关路由（无需认证）
		auth := v1.Group("/auth")
		{
			auth.POST("/wechat-login", authHandler.WechatLogin)
			auth.GET("/app-version", authHandler.GetAppVersion)
		}

		// 需要认证的路由
		authRequired := v1.Group("")
		authRequired.Use(middleware.Auth(cfg))
		{
			// 认证相关（需要token）
			authRequired.POST("/auth/refresh-token", authHandler.RefreshToken)
			authRequired.GET("/auth/user-info", authHandler.GetUserInfo)
			authRequired.PUT("/auth/user-info", authHandler.UpdateUserInfo)

			// 文件上传
			authRequired.POST("/upload", uploadHandler.Upload)
		}
	}

	return r
}
