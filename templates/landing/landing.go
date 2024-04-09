package templates

import t "github.com/rirachii/golivechat/templates"

var LandingPage = t.TemplateData{
	TemplateName: "landing",
}

type TemplateLandingPage struct {
	Title string
}
