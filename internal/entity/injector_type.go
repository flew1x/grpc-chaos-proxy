package entity

type InjectorType string

const (
	AbortType      InjectorType = "abort"
	ChaosType      InjectorType = "chaos"
	DelayType      InjectorType = "delay"
	SpammerType    InjectorType = "spammer"
	NetworkType    InjectorType = "network"
	HeaderType     InjectorType = "header"
	RateLimitType  InjectorType = "rate_limit"
	DisconnectType InjectorType = "disconnect"
	CodeType       InjectorType = "code"
	ScriptType     InjectorType = "script"
)

func (i InjectorType) String() string {
	return string(i)
}
