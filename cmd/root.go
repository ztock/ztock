package cmd

import (
	"context"
	"fmt"
	"os"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/ztock/ztock/internal/config"
	"github.com/ztock/ztock/pkg/stock"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	name    = "ztock"
	version = "1.0.0"
)

var cfg *config.Config

var cfgFile string

var rootCmd = &cobra.Command{
	Use:     name,
	Version: version,
	Short:   "Show stock real-time data tools",
	Long: `A command line tool to display real-time stock 
information and analysis results.
Complete documentation is available at https://github.com/ztock/ztock`,
	SilenceUsage: true,
	Args:         cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		// Init logger
		initLog(cfg)

		// Parse args
		cfg.Number = args[0]
		logrus.Debugf("Load config success: %#v", cfg)

		// Get stock info
		s := stock.NewStockContext(ctx, cfg.Platform, cfg)
		data, err := s.Get()
		logrus.Debugf("Get stock data success: %#v", data)
		if err != nil {
			return err
		}

		// Print data
		prettyPrint(data)
		return nil
	},
}

// Execute is the entry point of the command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		logrus.Debugf("Execute error: %#v", err)
	}
}

func init() {
	// Init default config
	var err error
	cfg, err = config.New()
	if err != nil {
		panic(err)
	}

	// Initialize cobra
	cobra.OnInitialize(initConfig)

	// Add flags
	flagSet := rootCmd.PersistentFlags()
	flagSet.VarP(&cfg.Platform, "platform", "p", "set the source platform for stock data")
	flagSet.VarP(&cfg.Index, "index", "i", "set the stock market index")
	flagSet.StringVar(&cfg.LogLevel, "log-level", cfg.LogLevel, "set the level that is used for logging")
	flagSet.StringVar(&cfg.LogFormat, "log-format", cfg.LogFormat, "set the format that is used for logging")
	flagSet.StringVar(&cfgFile, "config", "", "config file (default is $HOME/.ztock/config.yaml)")

	if err := viper.BindPFlags(rootCmd.PersistentFlags()); err != nil {
		panic(err)
	}
}

// initConfig reads in config file and ENV variables if set
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		cfgPath := fmt.Sprintf("%s/.%s", home, name)
		viper.AddConfigPath(cfgPath)
	}

	viper.SetEnvPrefix(name)
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := viper.ReadInConfig(); err == nil {
		logrus.Debugf("Using config file: %s", viper.ConfigFileUsed())
	}

	// Unmarshal config
	if err := viper.Unmarshal(&cfg); err != nil {
		logrus.Fatalf(errors.Wrap(err, "cannot unmarshal config").Error())
	}
}

func initLog(cfg *config.Config) {
	// Reset log format
	if cfg.LogFormat == "json" {
		logrus.SetFormatter(&logrus.JSONFormatter{})
	}

	// Set the configured log level
	if level, err := logrus.ParseLevel(cfg.LogLevel); err == nil {
		logrus.SetLevel(level)
	}
}

func prettyPrint(s *stock.Stock) {
	timeLayout := "2006-01-02 15:04:05"

	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)

	// Set table style
	if s.PercentageChange[0:1] == "-" {
		t.SetStyle(table.StyleColoredBlackOnGreenWhite)
	} else {
		t.SetStyle(table.StyleColoredBlackOnMagentaWhite)
	}

	// Set table header
	t.AppendHeader(table.Row{"Number", "Current Price", "Percentage Change", "Opening Price", "Previous Closing Price", "High Price", "Low Price", "Date"})
	t.AppendRows([]table.Row{
		{s.Number, s.CurrentPrice, s.PercentageChange, s.OpeningPrice, s.PreviousClosingPrice, s.HighPrice, s.LowPrice, s.Date.Format(timeLayout)},
	})
	t.AppendSeparator()

	t.Render()
}
