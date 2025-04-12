package domain

type IUtils interface {
	TemplateToString(path string, data any) (string, error)
	ReadFile(name string) ([]byte, error)
}
