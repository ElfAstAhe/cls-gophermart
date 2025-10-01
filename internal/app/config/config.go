// Package config
/*
  Iteration 5

  Configuration params priority :

  1 - ENV vars
  2 - CLI params
  3 - Default values
*/
package config

import (
	"flag"
	"fmt"
	"os"
	"sync"

	"github.com/caarlos0/env/v6"
	"go.uber.org/zap"
)

type Config struct {
	onceRun        sync.Once
	AppName        string      `json:"app_name,omitempty"`
	ProjectStage   string      `json:"project_stage,omitempty" env:"PROJECT_STAGE" envDefault:"DEV"`
	LogLevel       string      `json:"log_level,omitempty" env:"LOG_LEVEL" envDefault:"INFO"`
	LogFilePath    string      `json:"log_file_path,omitempty" env:"LOG_FILE_PATH"`
	BaseURL        string      `json:"base_url,omitempty" env:"BASE_URL"`
	HTTP           *HTTPConfig `json:"http,omitempty" env:"RUN_ADDRESS"`
	DBDsn          string      `json:"db_dsn,omitempty" env:"DATABASE_URI"`
	AccrualBaseURI string      `json:"accrual_base_uri,omitempty" env:"ACCRUAL_SYSTEM_ADDRESS"`
}

// Flags
const (
	FlagAppName        string = "p"
	FlagProjectStage   string = "s"
	FlagLogLevel       string = "l"
	FlagLogFilePath    string = "lf"
	FlagBaseURL        string = "b"
	FlagHTTPInterface  string = "a"
	FlagDBInterface    string = "d"
	FLagAccrualBaseURI string = "r"
)

// Environment variables
const (
	EnvProjectStage   string = "PROJECT_STAGE"
	EnvLogLevel       string = "LOG_LEVEL"
	EnvLogFilePath    string = "LOG_FILE_PATH"
	EnvBaseURL        string = "BASE_URL"
	EnvHTTPInterface  string = "RUN_ADDRESS"
	EnvDatabaseDSN    string = "DATABASE_URI"
	EnvAccrualBaseURI string = "ACCRUAL_SYSTEM_ADDRESS"
)

func NewConfig() *Config {
	var cfg = defaultConfig()

	cfg.onceRun.Do(cfg.initFlags)

	return cfg
}

func defaultConfig() *Config {
	return &Config{
		AppName:      DefaultAppName,
		ProjectStage: DefaultStage,
		LogLevel:     DefaultLogLevel,
		LogFilePath:  DefaultLogFilePath,
	}
}

func (c *Config) LoadConfig() error {
	fmt.Println("Parse cli params")
	var err = c.loadCli()
	if err != nil {
		return err
	}

	fmt.Println("Parse env params")
	err = c.loadEnv()
	if err != nil {
		return err
	}

	fmt.Printf("Final config:  [%+v]\n", c)

	return nil
}

func (c *Config) loadCli() error {
	flag.Parse()

	fmt.Printf("Config: %+v\r\n", c)

	return nil
}

func (c *Config) loadEnv() error {
	err := env.Parse(c)
	if err != nil {
		return err
	}

	err = parseFlag(EnvHTTPInterface, c.HTTP)
	if err != nil {
		return err
	}

	fmt.Printf("Config: %+v\r\n", c)

	return nil
}

func parseFlag(env string, value flag.Value) error {
	var envVar = os.Getenv(env)
	if envVar == "" {
		return nil
	}

	fmt.Printf("[DEBUG] Config: ENV [%s] VALUE [%+v]\r\n", env, value)

	err := value.Set(envVar)
	if err != nil {
		return err
	}

	return nil
}

func (c *Config) initFlags() {
	flag.StringVar(&c.AppName, FlagAppName, DefaultAppName, "application name")
	flag.StringVar(&c.ProjectStage, FlagProjectStage, ProjectStageDevelopment, "project stage")
	flag.StringVar(&c.LogLevel, FlagLogLevel, zap.InfoLevel.CapitalString(), "log level")
	flag.StringVar(&c.LogFilePath, FlagLogFilePath, DefaultLogFilePath, "log file path")
	flag.StringVar(&c.BaseURL, FlagBaseURL, DefaultBaseURL, "base url")
	flag.Var(c.HTTP, FlagHTTPInterface, "http interface")
	flag.StringVar(&c.DBDsn, FlagDBInterface, DefaultDBDsn, "database dsn")
}
