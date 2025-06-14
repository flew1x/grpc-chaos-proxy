package config

type Match struct {
	Service     string `mapstructure:"service" yaml:"service"`           // Service name to match
	MethodRegex string `mapstructure:"method_regex" yaml:"method_regex"` // Regular expression to match method names
}

type Rule struct {
	Name     string `mapstructure:"name" yaml:"name"`         // Name of the rule
	Match    Match  `mapstructure:"match" yaml:"match"`       // Match criteria for the rule
	Action   Action `mapstructure:"action" yaml:"action"`     // Action to take when the rule matches
	Disabled bool   `mapstructure:"disabled" yaml:"disabled"` // Whether the rule is disabled
}

type Listener struct {
	Address string `mapstructure:"address" yaml:"address"` // Address to listen on
}

type Backend struct {
	Address string `mapstructure:"address" yaml:"address"` // Backend gRPC service address
}

type Config struct {
	Listener Listener `mapstructure:"listener" yaml:"listener"`
	Backend  Backend  `mapstructure:"backend" yaml:"backend"`
	Rules    []Rule   `mapstructure:"rules" yaml:"rules"`
}
