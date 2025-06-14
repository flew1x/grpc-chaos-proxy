package entity

type InjectorType string

const (
	AbortType   InjectorType = "abort"
	ChaosType   InjectorType = "chaos"
	DelayType   InjectorType = "delay"
	SpammerType InjectorType = "spammer"
)

func (i InjectorType) String() string {
	return string(i)
}
