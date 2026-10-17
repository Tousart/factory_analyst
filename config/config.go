package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
)

const (
	envConfigPath = "CONFIG_PATH"

	envFactoryAnalystStgHost  = "FACTORY_ANALYST_STG_HOST"
	envFactoryAnalystProdHost = "FACTORY_ANALYST_PROD_HOST"

	envFactoryAnalystStgGatewayHost  = "FACTORY_ANALYST_STG_GATEWAY_HOST"
	envFactoryAnalystProdGatewayHost = "FACTORY_ANALYST_PROD_GATEWAY_HOST"
)

type Config struct {
	Environment Environment `yaml:"environment"`
	Hosts       Hosts       `yaml:"hosts"`
}

// Environment определяет, в каком окружении запускается приложение.
type Environment string

const (
	EnvironmentLocal Environment = "local"
	EnvironmentStg   Environment = "stg"
	EnvironmentProd  Environment = "prod"
)

func (e Environment) Valid() bool {
	switch e {
	case EnvironmentLocal, EnvironmentStg, EnvironmentProd:
		return true
	default:
		return false
	}
}

type Hosts struct {
	FactoryAnalyst `yaml:"factory_analyst"`
}

// FactoryAnalyst хранит адреса gRPC-сервера и grpc-gateway.
type FactoryAnalyst struct {
	LocalHost string `yaml:"local_host"`
	StgHost   string `yaml:"-"`
	ProdHost  string `yaml:"-"`

	LocalGatewayHost string `yaml:"local_gateway_host"`
	StgGatewayHost   string `yaml:"-"`
	ProdGatewayHost  string `yaml:"-"`
}

// Host возвращает адрес gRPC-сервера для текущего окружения.
func (c *Config) Host() string {
	switch c.Environment {
	case EnvironmentStg:
		return c.Hosts.FactoryAnalyst.StgHost
	case EnvironmentProd:
		return c.Hosts.FactoryAnalyst.ProdHost
	default:
		return c.Hosts.FactoryAnalyst.LocalHost
	}
}

// GatewayHost возвращает адрес grpc-gateway для текущего окружения.
func (c *Config) GatewayHost() string {
	switch c.Environment {
	case EnvironmentStg:
		return c.Hosts.FactoryAnalyst.StgGatewayHost
	case EnvironmentProd:
		return c.Hosts.FactoryAnalyst.ProdGatewayHost
	default:
		return c.Hosts.FactoryAnalyst.LocalGatewayHost
	}
}

// LoadFromEnv читает .env (если он есть) и загружает конфиг по пути из CONFIG_PATH.
func LoadFromEnv() (*Config, error) {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		return nil, fmt.Errorf("load .env: %w", err)
	}

	path := os.Getenv(envConfigPath)
	if path == "" {
		return nil, fmt.Errorf("env %s is not set", envConfigPath)
	}

	return Load(path)
}

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %s: %w", path, err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config %s: %w", path, err)
	}

	if !cfg.Environment.Valid() {
		return nil, fmt.Errorf("config %s: unknown environment %q, want %q, %q or %q",
			path, cfg.Environment, EnvironmentLocal, EnvironmentStg, EnvironmentProd)
	}

	cfg.Hosts.FactoryAnalyst.StgHost = os.Getenv(envFactoryAnalystStgHost)
	cfg.Hosts.FactoryAnalyst.ProdHost = os.Getenv(envFactoryAnalystProdHost)
	cfg.Hosts.FactoryAnalyst.StgGatewayHost = os.Getenv(envFactoryAnalystStgGatewayHost)
	cfg.Hosts.FactoryAnalyst.ProdGatewayHost = os.Getenv(envFactoryAnalystProdGatewayHost)

	if err := cfg.validateHost(path); err != nil {
		return nil, err
	}

	if err := cfg.validateGatewayHost(path); err != nil {
		return nil, err
	}

	return &cfg, nil
}

// validateHost проверяет, что задан адрес gRPC-сервера для выбранного окружения.
func (c *Config) validateHost(path string) error {
	if c.Host() != "" {
		return nil
	}

	switch c.Environment {
	case EnvironmentStg:
		return fmt.Errorf("env %s is not set", envFactoryAnalystStgHost)
	case EnvironmentProd:
		return fmt.Errorf("env %s is not set", envFactoryAnalystProdHost)
	default:
		return fmt.Errorf("config %s: hosts.factory_analyst.local_host is empty", path)
	}
}

// validateGatewayHost проверяет, что задан адрес grpc-gateway для выбранного окружения.
func (c *Config) validateGatewayHost(path string) error {
	if c.GatewayHost() != "" {
		return nil
	}

	switch c.Environment {
	case EnvironmentStg:
		return fmt.Errorf("env %s is not set", envFactoryAnalystStgGatewayHost)
	case EnvironmentProd:
		return fmt.Errorf("env %s is not set", envFactoryAnalystProdGatewayHost)
	default:
		return fmt.Errorf("config %s: hosts.factory_analyst.local_gateway_host is empty", path)
	}
}
