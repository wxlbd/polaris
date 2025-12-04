//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/google/wire"
	"github.com/wxlbd/polaris/internal/application/service"
	"github.com/wxlbd/polaris/internal/infrastructure/config"
	"github.com/wxlbd/polaris/internal/infrastructure/logger"
	"github.com/wxlbd/polaris/internal/infrastructure/persistence"
	"github.com/wxlbd/polaris/internal/infrastructure/wechat"
	"github.com/wxlbd/polaris/internal/interface/http/handler"
	"github.com/wxlbd/polaris/internal/interface/http/router"
)

// InitApp 初始化应用(Wire自动生成)
func InitApp(cfg *config.Config) (*App, error) {
	wire.Build(
		// 基础设施层
		logger.NewLogger, // 日志系统
		persistence.NewDatabase,
		persistence.NewRedis, // Redis 客户端
		wechat.NewClient,     // 微信 SDK 客户端

		// 仓储层
		persistence.NewUserRepository,
		persistence.NewAppVersionRepository, // 应用版本仓储

		// 应用服务层
		service.NewAuthService,
		service.NewUploadService,     // 文件上传服务
		service.NewAppVersionService, // 应用版本服务

		// HTTP处理器
		handler.NewAuthHandler,
		handler.NewUploadHandler, // 文件上传处理器

		// 路由
		router.NewRouter,

		// 应用
		NewApp,
	)
	return &App{}, nil
}
