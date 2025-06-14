package config

type DelayAction struct {
	MinMS int `mapstructure:"min_ms" validate:"gte=0"` // Minimum delay in milliseconds
	MaxMS int `mapstructure:"max_ms" validate:"gte=0"` // Maximum delay in milliseconds
}

type AbortAction struct {
	Code       string `mapstructure:"code"`       // gRPC code to return
	Percentage int    `mapstructure:"percentage"` // Percentage of requests to abort
}

type SpammerAction struct {
	Count       int          `mapstructure:"count" yaml:"count"` // Number of spam requests to send
	DelayAction *DelayAction `mapstructure:"delay" yaml:"delay"` // Optional delay between spam requests
}

type ChaosAction struct {
	Actions []Action `mapstructure:"actions" yaml:"actions"`
}

type Action struct {
	Delay   *DelayAction   `mapstructure:"delay" yaml:"delay"`     // Optional delay action
	Abort   *AbortAction   `mapstructure:"abort" yaml:"abort"`     // Whether to abort the request
	Chaos   *ChaosAction   `mapstructure:"chaos" yaml:"chaos"`     // Chaos actions to apply
	Spammer *SpammerAction `mapstructure:"spammer" yaml:"spammer"` // Spammer action to send multiple requests
}
