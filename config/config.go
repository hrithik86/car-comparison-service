package config

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"strconv"
	"time"
)

var config *Configuration

type Configuration struct {
	module       string
	port         int64
	profile      string
	logLevel     string
	logLocation  string
	enableSqlLog bool
	enableSentry bool
	sentryDsn    string
	dbConfig     DbConfig
	redisConfig  RedisConfig
	queryTimeout int64
}

type RedisConfig struct {
	Host                  string
	User                  string
	Password              string
	ReadTimeout           time.Duration
	WriteTimeout          time.Duration
	RedisKeyExpiryTimeout time.Duration
}

type DbConfig struct {
	Host            string
	Port            int64
	Name            string
	User            string
	Password        string
	MaxIdleConn     int64
	ConnMaxIdleTime int64
	ConnMaxLifeTime int64
	MaxOpenConn     int64
}

func Load(commandArgs []string) {
	module := commandArgs[1]
	viper.SetDefault("log_level", "debug")
	viper.AutomaticEnv()

	viper.AddConfigPath("./profiles")
	viper.SetConfigType("yml")

	if len(commandArgs) > 2 {
		viper.SetConfigName(commandArgs[2])

		err := viper.ReadInConfig()
		if err != nil {
			fmt.Println("error while reading config file")
			panic(err)
		}

		fmt.Println("In memory config read successfully")
	} else {
		logrus.Fatal("Failed to startup server, config file name missing")
	}

	config = &Configuration{
		module:       module,
		port:         getIntOrPanic("port"),
		profile:      getStringOrPanic("profile"),
		logLevel:     getStringOrPanic("log_level"),
		logLocation:  getStringOrPanic("log_location"),
		enableSqlLog: getBoolOrPanic("enable_sql_log"),
		queryTimeout: getIntOrPanic("query_timeout"),
		dbConfig: DbConfig{
			Host:            getStringOrPanic("db_host"),
			Port:            getIntOrPanic("db_port"),
			Name:            getStringOrPanic("db_name"),
			User:            getStringOrPanic("db_user"),
			Password:        getStringOrPanic("db_password"),
			MaxIdleConn:     getIntOrPanic("db_max_idle_conn"),
			ConnMaxIdleTime: getIntOrPanic("db_conn_max_idle_time"),
			ConnMaxLifeTime: getIntOrPanic("db_conn_max_life_time"),
			MaxOpenConn:     getIntOrPanic("db_max_open_conn"),
		},
		redisConfig: RedisConfig{
			Host:                  getStringOrPanic("redis_host"),
			User:                  getStringOrPanic("redis_user"),
			Password:              getStringOrPanic("redis_password"),
			ReadTimeout:           time.Duration(getIntOrPanic("redis_read_connection_timeout_milliseconds")) * time.Millisecond,
			WriteTimeout:          time.Duration(getIntOrPanic("redis_write_connection_timeout_milliseconds")) * time.Millisecond,
			RedisKeyExpiryTimeout: time.Duration(getIntOrPanic("redis_key_expiry_timeout_milliseconds")) * time.Millisecond,
		},
	}

}

func getIntOrPanic(key string) int64 {
	checkKey(key)
	v, err := strconv.Atoi(viper.GetString(key))
	panicIfErrorForKey(err, key)
	return int64(v)
}

func getStringOrPanic(key string) string {
	checkKey(key)
	return viper.GetString(key)
}

func getBoolOrPanic(key string) bool {
	if !viper.IsSet(key) {
		return false
	}

	v, err := strconv.ParseBool(viper.GetString(key))
	panicIfErrorForKey(err, key)
	return v
}

func checkKey(key string) {
	if !viper.IsSet(key) {
		panicIfError(fmt.Errorf("%s key is not set", key))
	}
}

func panicIfErrorForKey(err error, key string) {
	if err != nil {
		panicIfError(fmt.Errorf("could not parse key: %s. Error: %v", key, err))
	}
}

func panicIfError(err error) {
	if err != nil {
		panic(fmt.Errorf("unable to load config: %v", err))
	}
}

func DbConf() DbConfig {
	return config.dbConfig
}

func QueryTimeout() time.Duration {
	return time.Duration(config.queryTimeout) * time.Millisecond
}

func Port() int64 {
	return config.port
}

func LogLevel() string {
	return config.logLevel
}

func RedisConf() RedisConfig {
	return config.redisConfig
}
