package config

import (
	"fmt"

	"github.com/kelseyhightower/envconfig"
)

type config struct {
	Port                     string `envconfig:"PORT" default:"8080"`
	ExtractorConcurrentLimit int    `envconfig:"EXTRACTOR_CONCURRENT_LIMIT" default:"2"`
	RateLimitRPM             int    `envconfig:"RATE_LIMIT_RPM" default:"60"`
}

func GetConfig() (*config, error) {
	var config config
	err := envconfig.Process("", &config)
	if err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}

	return &config, nil
}
