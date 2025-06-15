package config

type DelayAction struct {
	MinMS int `mapstructure:"min_ms" validate:"gte=0"` // Minimum delay in milliseconds
	MaxMS int `mapstructure:"max_ms" validate:"gte=0"` // Maximum delay in milliseconds
}

type AbortAction struct {
	Code       string `mapstructure:"code"`       // gRPC code to return
	Percentage int    `mapstructure:"percentage"` // Percentage of requests to abort
}

type DisconnectAction struct {
	Percentage int `mapstructure:"percentage"` // Percentage of requests to disconnect
}

type SpammerAction struct {
	Count       int          `mapstructure:"count" yaml:"count"` // Number of spam requests to send
	DelayAction *DelayAction `mapstructure:"delay" yaml:"delay"` // Optional delay between spam requests
}

type ChaosAction struct {
	Actions []Action `mapstructure:"actions" yaml:"actions"`
}

type HeaderAction struct {
	Headers map[string]string `mapstructure:"headers" yaml:"headers"` // Headers to add to the request
}

type NetworkAction struct {
	LossPercentage      int `mapstructure:"loss_percentage" yaml:"loss_percentage"`           // % loss of packets
	DuplicatePercentage int `mapstructure:"duplicate_percentage" yaml:"duplicate_percentage"` // % duplicated packets
	ReorderPercentage   int `mapstructure:"reorder_percentage" yaml:"reorder_percentage"`     // % change of packet reordering
	ThrottleMS          int `mapstructure:"throttle_ms" yaml:"throttle_ms"`                   // delay in milliseconds to throttle packets
}

type RateLimiterAction struct {
	BurstSize int `mapstructure:"burst_size" yaml:"burst_size"` // Maximum burst size for rate limiting
	RateLimit int `mapstructure:"rate_limit" yaml:"rate_limit"` // Rate limit in requests per second
}

type CodeAction struct {
	Code        string            `mapstructure:"code" yaml:"code"`                       // gRPC code to return
	Message     string            `mapstructure:"message" yaml:"message"`                 // Custom error message
	Percentage  int               `mapstructure:"percentage" yaml:"percentage"`           // Probability to inject error (0-100)
	Metadata    map[string]string `mapstructure:"metadata" yaml:"metadata"`               // Metadata to add to response
	DelayMS     int               `mapstructure:"delay_ms" yaml:"delay_ms"`               // Delay before returning error (ms)
	OnlyOn      []string          `mapstructure:"only_on_methods" yaml:"only_on_methods"` // Methods to apply to
	RepeatCount int               `mapstructure:"repeat_count" yaml:"repeat_count"`       // How many times to repeat error
}

type ScriptAction struct {
	Language   string            `mapstructure:"language" yaml:"language"`     // Language of the script (e.g., "sh", "bash")
	Source     string            `mapstructure:"source" yaml:"source"`         // Source code of the script
	Entrypoint string            `mapstructure:"entrypoint" yaml:"entrypoint"` // Name of the entrypoint function to call
	TimeoutMS  int               `mapstructure:"timeout_ms" yaml:"timeout_ms"` // Timeout for script execution in milliseconds
	Env        map[string]string `mapstructure:"env" yaml:"env"`               // Variables to set in the script environment
	Args       []string          `mapstructure:"args" yaml:"args"`             // Arguments to pass to the script
}

type Action struct {
	Delay       *DelayAction       `mapstructure:"delay" yaml:"delay"`               // Delay action
	Abort       *AbortAction       `mapstructure:"abort" yaml:"abort"`               // Abort whether to abort the request
	Chaos       *ChaosAction       `mapstructure:"chaos" yaml:"chaos"`               // Chaos actions to apply
	Spammer     *SpammerAction     `mapstructure:"spammer" yaml:"spammer"`           // Spammer action to send multiple requests
	Network     *NetworkAction     `mapstructure:"network" yaml:"network"`           // Network action
	Header      *HeaderAction      `mapstructure:"header" yaml:"header"`             // Header action to add headers
	RateLimiter *RateLimiterAction `mapstructure:"rate_limiter" yaml:"rate_limiter"` // Rate limiter action
	Disconnect  *DisconnectAction  `mapstructure:"disconnect" yaml:"disconnect"`     // Disconnect action to drop connections
	Code        *CodeAction        `mapstructure:"code" yaml:"code"`                 // Code action to return specific gRPC code
	Script      *ScriptAction      `mapstructure:"script" yaml:"script"`             // Script action to execute custom logic
}
