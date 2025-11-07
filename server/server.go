package server

import (
	httpapi "oracle.com/oracle/my-go-oracle-app/api/http"
	config "oracle.com/oracle/my-go-oracle-app/configs"
)

// Init to initiate all DI for service handler implementation
func InitHttp(config *config.Config) error {

	httpserver := httpapi.Server{
		Cfg: config,
	}

	return runHTTPServer(httpserver, config.ServerHttpPort)
}
