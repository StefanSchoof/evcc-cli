package cmd

import (
	"fmt"
	"strings"
	"time"

	"evcc-cli/cmd/loadpoint"
	"evcc-cli/internal/config"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfg config.Config

var rootCmd = &cobra.Command{
	Use:           "evcc-cli",
	Short:         "CLI for evcc",
	SilenceUsage:  true,
	SilenceErrors: true,
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		cfg.Host = strings.TrimRight(viper.GetString("host"), "/")
		cfg.Timeout = viper.GetDuration("timeout")
		cfg.Raw = viper.GetBool("raw")
		cfg.Insecure = viper.GetBool("insecure")
		return cfg.Validate()
	},
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	cobra.OnInitialize(initConfig)

	rootCmd.PersistentFlags().String("host", "http://localhost:7070", "evcc base URL")
	rootCmd.PersistentFlags().Duration("timeout", 10*time.Second, "HTTP timeout")
	rootCmd.PersistentFlags().Bool("raw", false, "print original API response body when available")
	rootCmd.PersistentFlags().Bool("insecure", false, "skip TLS certificate verification")

	mustBind("host")
	mustBind("timeout")
	mustBind("raw")
	mustBind("insecure")

	viper.SetEnvPrefix("EVCC")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	rootCmd.AddCommand(newVersionCmd())
	rootCmd.AddCommand(newCompletionCmd())
	rootCmd.AddCommand(newStateCmd())
	rootCmd.AddCommand(loadpoint.NewCmd(newGeneratedClient, func() bool { return cfg.Raw }))
}

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("$HOME/.config/evcc-cli")
	_ = viper.ReadInConfig()
}

func mustBind(flagName string) {
	if err := viper.BindPFlag(flagName, rootCmd.PersistentFlags().Lookup(flagName)); err != nil {
		panic(fmt.Sprintf("failed to bind flag %s: %v", flagName, err))
	}
}
