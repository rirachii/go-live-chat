package templates

import (
	t "github.com/rirachii/golivechat/templates"
)

// "location-data" template
var LocationDataFragment = t.TemplateData{
	TemplateName: "location-data",
}

type FragmentLocationData struct {
	County  string
	City    string
	Country string
}
