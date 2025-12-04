package wire

import (
	"github.com/gin-gonic/gin"
	"github.com/wxlbd/polaris/internal/infrastructure/config"
)

// App 应用程序
type App struct {
	Config *config.Config
	Router *gin.Engine
}

// NewApp 创建应用实例
func NewApp(
	cfg *config.Config,
	router *gin.Engine,
) *App {
	return &App{
		Config: cfg,
		Router: router,
	}
}
