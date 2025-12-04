package main

import (
	"flag"
	"log"
	"os"

	"github.com/wxlbd/polaris/internal/domain/entity"
	"github.com/wxlbd/polaris/internal/infrastructure/config"
	"github.com/wxlbd/polaris/internal/infrastructure/logger"
	"github.com/wxlbd/polaris/internal/infrastructure/persistence"
)

func main() {
	var command string
	flag.StringVar(&command, "command", "up", "Migration command: up or down")
	flag.Parse()

	if len(os.Args) > 1 {
		command = os.Args[1]
	}

	// 加载配置
	cfg, err := config.Load("config/config.yaml")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	if err := logger.Init(cfg.Log); err != nil {
		log.Fatalf("Failed to init logger: %v", err)
	}

	// 连接数据库
	db, err := persistence.NewDatabase(cfg)
	if err != nil {
		log.Fatalf("Failed to connect database: %v", err)
	}

	// 自动迁移
	// 注意：这里使用 GORM 的 AutoMigrate 功能，它会自动创建表
	// 对于生产环境，建议使用 golang-migrate 等工具进行版本控制
	if command == "up" {
		log.Println("Running migrations up...")
		err = db.AutoMigrate(
			&entity.User{},
			&entity.AppVersion{},
		)
		if err != nil {
			log.Fatalf("Migration failed: %v", err)
		}
		log.Println("Migration completed successfully")
	} else if command == "down" {
		log.Println("Migration down is not supported with AutoMigrate. Please drop tables manually if needed.")
	} else {
		log.Fatalf("Unknown command: %s", command)
	}
}
