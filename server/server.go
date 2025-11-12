package server

import (
	"fmt"
	"log/slog"
	"os"

	"oracle.com/oracle/my-go-oracle-app/api"
	httpapi "oracle.com/oracle/my-go-oracle-app/api/http"
	config "oracle.com/oracle/my-go-oracle-app/configs"
	"oracle.com/oracle/my-go-oracle-app/infra/database/sql"
	"oracle.com/oracle/my-go-oracle-app/service"
	"oracle.com/oracle/my-go-oracle-app/service/member"
)

// Init to initiate all DI for service handler implementation
func InitHttp(config *config.Config) error {
	// baseRepo := getBaseRepository(config)
	baseRepo := getBaseRepository(config)

	memberRepo := member.NewMemberRepository(baseRepo)
	memberService := member.NewMemberService(memberRepo)

	httpserver := httpapi.Server{
		Cfg:           config,
		MemberService: memberService,
		HealthCheck: api.HealthChecker{
			Master: baseRepo.MasterDB,
			Slave:  baseRepo.SlaveDB,
		},
	}

	return runHTTPServer(httpserver, config.ServerHttpPort)
}

func getBaseRepository(config *config.Config) service.BaseRepository {

	var dbMasterURL, dbSlaveURL string
	if config.OracleLibDir == "" {
		dbMasterURL = fmt.Sprintf("%s/%s@%s:%s/%s", config.OracleMasterUsername, config.OracleMasterPassword, config.OracleMasterHost, config.OracleMasterPort, config.OracleMasterDatabase)
		dbSlaveURL = fmt.Sprintf("%s/%s@%s:%s/%s", config.OracleSlaveUsername, config.OracleSlavePassword, config.OracleSlaveHost, config.OracleSlavePort, config.OracleSlaveDatabase)
	} else {
		masterConn := fmt.Sprintf("%s:%s/%s", config.OracleMasterHost, config.OracleMasterPort, config.OracleMasterDatabase)
		slaveConn := fmt.Sprintf("%s:%s/%s", config.OracleSlaveHost, config.OracleSlavePort, config.OracleSlaveDatabase)
		dbMasterURL = fmt.Sprintf(`user="%s"
                                password="%s"
                                connectString="%s"
                                libDir="%s"`,
			config.OracleMasterUsername,
			config.OracleMasterPassword,
			masterConn,
			config.OracleLibDir,
		)
		dbSlaveURL = fmt.Sprintf(`user="%s"
                                password="%s"
                                connectString="%s"
                                libDir="%s"`,
			config.OracleSlaveUsername,
			config.OracleSlavePassword,
			slaveConn,
			config.OracleLibDir,
		)
	}
	//init db config

	masterDB, err := sql.OpenMasterDB("godror", dbMasterURL, config.OracleMaxOpenConnection, config.OracleMaxIdleConnection, config.OracleConnMaxIdleTime, config.OracleConnMaxLifeTime)
	if err != nil {
		slog.Error(fmt.Sprintf("init master DB failed: %v", err))
		os.Exit(1)
	}
	slaveDB, err := sql.OpenSlaveDB("godror", dbSlaveURL, config.OracleMaxOpenConnection, config.OracleMaxIdleConnection, config.OracleConnMaxIdleTime, config.OracleConnMaxLifeTime)
	if err != nil {
		slog.Error(fmt.Sprintf("init slave DB failed: %v", err))
		os.Exit(1)
	}

	return service.BaseRepository{
		MasterDB: masterDB,
		SlaveDB:  slaveDB,
	}
}
