package config

import (
	"fmt"
)

// Config holds all the runtime config information
type Config struct {
	// Stock number
	Number string `mapstructure:"number"`

	// Stock market index
	Index IndexType `mapstructure:"index"`

	// Source platform for stock data
	Platform PlatformType `mapstructure:"platform"`

	// LogLevel is the level with to log for this config
	LogLevel string `mapstructure:"log_level"`

	// LogFormat is the format that is used for logging
	LogFormat string `mapstructure:"log_format"`
}

// PlatformType is the type of
type PlatformType string

const (
	// Sina platform
	SinaPlatformType PlatformType = "sina"
)

func (p *PlatformType) String() string {
	return fmt.Sprint(*p)
}

func (p *PlatformType) Set(value string) error {
	*p = PlatformType(value)
	return nil
}

func (p *PlatformType) Type() string {
	return "platform"
}

type IndexType string

const (
	// ShangHai Index
	ShangHaiIndexType IndexType = "sh"

	// ShenZhen Index
	ShenZhenIndexType IndexType = "sz"
)

func (i *IndexType) String() string {
	return fmt.Sprint(*i)
}

func (i *IndexType) Set(value string) error {
	*i = IndexType(value)
	return nil
}

func (i *IndexType) Type() string {
	return "index"
}

const (
	// DefaultIndex is the default stock market index
	DefaultIndex = ShangHaiIndexType

	// DefaultPlatform is the default source platform
	DefaultPlatform = SinaPlatformType

	// DefaultLogLevel is the default logging level
	DefaultLogLevel = "warn"

	// DefaultLogFormat is the default format of the logger
	DefaultLogFormat = "text"

	// DefaultLogFormat is the default format of the logger
	DefaultConfigExt = "yaml"
)

// New returns a new Config
func New() (*Config, error) {
	return &Config{
		Index:     DefaultIndex,
		Platform:  DefaultPlatform,
		LogLevel:  DefaultLogLevel,
		LogFormat: DefaultLogFormat,
	}, nil
}
