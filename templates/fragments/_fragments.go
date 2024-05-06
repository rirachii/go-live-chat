package templates

import (
	t "github.com/rirachii/golivechat/templates"
)

// "chatroom" template
var LocationDataFragment = t.TemplateData{
	TemplateName: "location-data",
}

type FragmentLocationData struct {
	County   string
	City string
	Country string
}