package entity

type ScriptLanguage string

const (
	ScriptLanguageBash ScriptLanguage = "bash"
	ScriptLanguageSh   ScriptLanguage = "sh"
)

func (l ScriptLanguage) String() string {
	return string(l)
}
