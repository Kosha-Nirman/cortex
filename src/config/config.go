package config

import (
	"time"

	"github.com/spf13/viper"
)

type ResolverConfig struct {
	Timeout       time.Duration
	Threads       int
	DNSServers    []string
	WordlistPath  string
	EnableBrute   bool
	EnableCRT     bool
	EnableDNS     bool
	EnablePassive bool
	HTTPTimeout   time.Duration
}

type OutputConfig struct {
	Directory string `mapstructure:"directory"`
	Format    string `mapstructure:"format"`
}

type LoggingConfig struct {
	Level  string `mapstructure:"level"`
	Format string `mapstructure:"format"`
}

type Config struct {
	Resolver *ResolverConfig `mapstructure:"resolver"`
	Output   OutputConfig    `mapstructure:"output"`
	Logging  LoggingConfig   `mapstructure:"logging"`
}

func setDefaultConfig() {
	viper.SetDefault("resolver.timeout", 5*time.Second)
	viper.SetDefault("resolver.threads", 100)
	viper.SetDefault("resolver.dns_servers", []string{"8.8.8.8:53", "1.1.1.1:53", "9.9.9.9:53"})
	viper.SetDefault("resolver.wordlist_path", "")
	viper.SetDefault("resolver.enable_crt", true)
	viper.SetDefault("resolver.enable_dns", true)
	viper.SetDefault("resolver.enable_brute", true)
	viper.SetDefault("resolver.enable_passive", true)
	viper.SetDefault("resolver.http_timeout", 10*time.Second)

	viper.SetDefault("output.directory", "")
	viper.SetDefault("output.format", "markdown")

	viper.SetDefault("logging.level", "info")
	viper.SetDefault("logging.format", "text")
}

func GetDefaultResolverConfig() *ResolverConfig {
	return &ResolverConfig{
		Timeout:       5 * time.Second,
		Threads:       100,
		DNSServers:    []string{"8.8.8.8:53", "1.1.1.1:53", "9.9.9.9:53"},
		WordlistPath:  "",
		EnableCRT:     true,
		EnableDNS:     true,
		EnableBrute:   true,
		EnablePassive: true,
		HTTPTimeout:   10 * time.Second,
	}
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("cortex")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.cortex")
	viper.AddConfigPath("/etc/cortex/")

	setDefaultConfig()

	if err := viper.ReadInConfig(); err != nil {
		// If the config file is not found, we can proceed with defaults
		if _, ok := err.(viper.ConfigFileNotFoundError); !ok {
			return nil, err
		}
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, err
	}

	// If Resolver is nil (not set in config), use default resolver config
	if config.Resolver == nil {
		config.Resolver = GetDefaultResolverConfig()
	}

	return &config, nil
}
