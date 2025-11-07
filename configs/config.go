package config

import (
	"github.com/spf13/viper"
)

var (
	config *Config
)

// option defines configuration option
type option struct {
	configFolder string
	configFile   string
	configType   string
}

// Init initializes `config` from the default config file.
// use `WithConfigFile` to specify the location of the config file
func Init(opts ...Option) error {
	opt := &option{
		configFolder: getDefaultConfigFolder(),
		configFile:   getDefaultConfigFile(),
		configType:   getDefaultConfigType(),
	}

	for _, optFunc := range opts {
		optFunc(opt)
	}

	// Config File Path
	viper.AddConfigPath(opt.configFolder)
	// Config File Name
	viper.SetConfigName(opt.configFile)
	// Config File Type
	viper.SetConfigType(opt.configType)
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return err
	}

	// Reading variables without using the model for Kube SecretKey
	config = new(Config)
	config.OracleMasterHost = viper.GetString("ORACLE_MASTER_HOST")
	config.OracleMasterPort = viper.GetString("ORACLE_MASTER_PORT")
	config.OracleMasterDatabase = viper.GetString("ORACLE_MASTER_DATABASE")
	config.OracleMasterUsername = viper.GetString("ORACLE_MASTER_USERNAME")
	config.OracleMasterPassword = viper.GetString("ORACLE_MASTER_PASSWORD")
	config.OracleSlaveHost = viper.GetString("ORACLE_SLAVE_HOST")
	config.OracleSlavePort = viper.GetString("ORACLE_SLAVE_PORT")
	config.OracleSlaveDatabase = viper.GetString("ORACLE_SLAVE_DATABASE")
	config.OracleSlaveUsername = viper.GetString("ORACLE_SLAVE_USERNAME")
	config.OracleSlavePassword = viper.GetString("ORACLE_SLAVE_PASSWORD")

	//set default value for all config
	setDefault()

	err = viper.Unmarshal(&config)
	if err != nil {
		return err
	}

	return config.postprocess()
}

// Option define an option for config package
type Option func(*option)

// WithConfigFolder set `config` to use the given config folder
func WithConfigFolder(configFolder string) Option {
	return func(opt *option) {
		opt.configFolder = configFolder
	}
}

// WithConfigFile set `config` to use the given config file
func WithConfigFile(configFile string) Option {
	return func(opt *option) {
		opt.configFile = configFile
	}
}

// WithConfigType set `config` to use the given config type
func WithConfigType(configType string) Option {
	return func(opt *option) {
		opt.configType = configType
	}
}

// getDefaultConfigFolder get default config folder.
func getDefaultConfigFolder() string {
	configPath := "./configs/"

	return configPath
}

// getDefaultConfigFile get default config file.
func getDefaultConfigFile() string {
	return "config"
}

// getDefaultConfigType get default config type.
func getDefaultConfigType() string {
	return "yaml"
}

// Get config
func Get() *Config {
	if config == nil {
		config = &Config{}
	}
	return config
}

func setDefault() {
	//redis default
	viper.SetDefault("REDIS_CACHE_ADDRESS", "redis1:6379,redis2:6379,redis3:6379")
}

// postprocess several config
func (c *Config) postprocess() error {
	return nil
}
