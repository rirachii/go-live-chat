package templates

import t "github.com/rirachii/golivechat/templates"

var RegisterPage = t.TemplateData{
	TemplateName: "register",
}

type TemplateRegisterPage struct {
	Title string
}