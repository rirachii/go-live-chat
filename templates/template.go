package templates

type TemplateData struct {
	TemplateName string
}

type DefinedTemplate interface {
	TemplateName() string
	PrepareData(interface{}) error
}
