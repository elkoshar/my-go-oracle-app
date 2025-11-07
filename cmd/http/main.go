package main

import (
	"fmt"
	"log/slog"
	"os"

	config "oracle.com/oracle/my-go-oracle-app/configs"
	"oracle.com/oracle/my-go-oracle-app/pkg/helpers"
	"oracle.com/oracle/my-go-oracle-app/pkg/logger"
	"oracle.com/oracle/my-go-oracle-app/pkg/panics"
	"oracle.com/oracle/my-go-oracle-app/server"
)

// @title MY GO ORACLE APP API
// @version 0.1
// @description This service is to handle My Go ORACLE App. For more detail, please visit https://github.com/
// @contact.name Oracle Team
// @contact.url https://github.com/orgs/xxxx/projects/1
// @contact.email oracle.team@mail.com
// @BasePath /my-go-oracle-app
func main() {

	var (
		cfg *config.Config
	)

	// init config
	err := config.Init(
		config.WithConfigFile("config"),
		config.WithConfigType("env"),
	)
	if err != nil {
		slog.Warn(fmt.Sprintf("failed to initialize config: %v", err))
		os.Exit(1)
	}
	cfg = config.Get()

	//init logging
	logger.InitLogger(cfg)

	// init send to Slack when panics
	panics.SetOptions(&panics.Options{
		Env: helpers.GetEnvString(),
	})

	// init all DI for service handler implementation
	if err := server.InitHttp(cfg); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

}
