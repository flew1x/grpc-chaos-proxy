package cli

import (
	"fmt"
	"os"

	"github.com/flew1x/grpc-chaos-proxy/internal/bootstrap"
	"go.uber.org/fx"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	yamlExt        = "yaml"
	defaultCfgPath = "./internal/configs/dev.yaml"
)

// NewCLI returns the root cobra.Command with subcommands:
//
//	proxy run    – start the proxy with the specified config
//	rule enable  – enable a rule by name (hot‑reload)
//	rule disable – disable a rule by name
//	show config  – print the current config (JSON)
func NewCLI() *cobra.Command {
	var cfgPath string

	root := &cobra.Command{
		Use:               "grpc-chaos-proxy",
		Short:             "Chaos‑engineering proxy for gRPC",
		PersistentPreRunE: requireConfig(&cfgPath),
	}

	root.PersistentFlags().StringVarP(&cfgPath, "config", "c", defaultCfgPath, "Path to YAML config")

	root.AddCommand(
		newRunCmd(&cfgPath),
		newRuleCmd(&cfgPath),
		newShowConfigCmd(&cfgPath),
	)

	return root
}

func requireConfig(cfgPath *string) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		if cfgPath == nil || *cfgPath == "" {
			return fmt.Errorf("--config required")
		}

		return nil
	}
}

func newRunCmd(cfgPath *string) *cobra.Command {
	return &cobra.Command{
		Use:   "run",
		Short: "Run proxy with supplied config file",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := ValidateConfigFile(*cfgPath); err != nil {
				return err
			}
			app := fxApp(*cfgPath)
			app.Run()

			return nil
		},
	}
}

func newRuleCmd(cfgPath *string) *cobra.Command {
	ruleCmd := &cobra.Command{Use: "rule", Short: "Manage rules in live config"}

	ruleCmd.AddCommand(
		&cobra.Command{
			Use:   "enable [name]",
			Short: "Enable rule by name in config file & hot‑reload",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return toggleRule(*cfgPath, args[0], true)
			},
		},
		&cobra.Command{
			Use:   "disable [name]",
			Short: "Disable rule by name in config file & hot‑reload",
			Args:  cobra.ExactArgs(1),
			RunE: func(cmd *cobra.Command, args []string) error {
				return toggleRule(*cfgPath, args[0], false)
			},
		},
	)

	return ruleCmd
}

func newShowConfigCmd(cfgPath *string) *cobra.Command {
	return &cobra.Command{
		Use:   "show config",
		Short: "Print current config (JSON)",
		RunE: func(cmd *cobra.Command, args []string) error {
			cfg, err := LoadConfig(*cfgPath)
			if err != nil {
				return err
			}

			return PrintConfigJSON(cfg)
		},
	}
}

func toggleRule(path, rule string, enable bool) error {
	cfg, err := LoadConfig(path)
	if err != nil {
		return err
	}

	found := false

	for i := range cfg.Rules {
		if cfg.Rules[i].Name == rule {
			cfg.Rules[i].Disabled = !enable
			found = true

			break
		}
	}

	if !found {
		return fmt.Errorf("rule %s not found", rule)
	}

	v := viper.New()
	v.SetConfigFile(path)
	v.SetConfigType(yamlExt)

	if err := v.ReadInConfig(); err != nil {
		return err
	}

	if err := v.MergeConfigMap(map[string]any{"rules": cfg.Rules}); err != nil {
		return err
	}

	return v.WriteConfig()
}

func fxApp(cfgPath string) *fx.App {
	_ = os.Setenv("CONFIG_PATH", cfgPath)

	return bootstrap.NewApp()
}
