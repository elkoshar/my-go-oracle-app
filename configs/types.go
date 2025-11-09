package config

import "time"

type (
	// Config will holds mapped key value for service configuration
	Config struct {
		AppVersion                    string        `mapstructure:"APP_VERSION"`
		ServerHttpPort                string        `mapstructure:"SERVER_HTTP_PORT"`
		LogLevel                      string        `mapstructure:"LOG_LEVEL"`
		LogFormat                     string        `mapstructure:"LOG_FORMAT"`
		HttpReadTimeout               time.Duration `mapstructure:"HTTP_READ_TIMEOUT"`
		HttpWriteTimeout              time.Duration `mapstructure:"HTTP_WRITE_TIMEOUT"`
		HttpInboundTimeout            time.Duration `mapstructure:"HTTP_INBOUND_TIMEOUT"`
		HTTPTimeout                   time.Duration `mapstructure:"HTTP_TIMEOUT"`
		HTTPDebug                     bool          `mapstructure:"HTTP_DEBUG"`
		HTTPMaxIdleConnections        int           `mapstructure:"HTTP_MAX_IDLE_CONNECTIONS"`
		HTTPMaxIdleConnectionsPerHost int           `mapstructure:"HTTP_MAX_IDLE_CONNECTIONS_PER_HOST"`
		HTTPIdleConnectionTimeout     time.Duration `mapstructure:"HTTP_IDLE_CONNECTION_TIMEOUT"`

		OracleMasterHost     string `mapstructure:"ORACLE_MASTER_HOST"`
		OracleMasterPort     string `mapstructure:"ORACLE_MASTER_PORT"`
		OracleMasterDatabase string `mapstructure:"ORACLE_MASTER_DATABASE"`
		OracleMasterUsername string `mapstructure:"ORACLE_MASTER_USERNAME"`
		OracleMasterPassword string `mapstructure:"ORACLE_MASTER_PASSWORD"`

		OracleSlaveHost     string `mapstructure:"ORACLE_SLAVE_HOST"`
		OracleSlavePort     string `mapstructure:"ORACLE_SLAVE_PORT"`
		OracleSlaveDatabase string `mapstructure:"ORACLE_SLAVE_DATABASE"`
		OracleSlaveUsername string `mapstructure:"ORACLE_SLAVE_USERNAME"`
		OracleSlavePassword string `mapstructure:"ORACLE_SLAVE_PASSWORD"`

		OracleMaxOpenConnection int           `mapstructure:"ORACLE_MAX_OPEN_CONNECTION"`
		OracleMaxIdleConnection int           `mapstructure:"ORACLE_MAX_IDLE_CONNECTION"`
		OracleConnMaxIdleTime   time.Duration `mapstructure:"ORACLE_CONN_MAX_IDLE_TIME"`
		OracleConnMaxLifeTime   time.Duration `mapstructure:"ORACLE_CONN_MAX_LIFE_TIME"`

		OracleLibDir string `mapstructure:"ORACLE_LIB_DIR"`
	}
)
