package controllers

import "github.com/eaciit/knot/knot.v1"

type Dashboard struct {
	*BaseController
}

func (c *Dashboard) Default(k *knot.WebContext) interface{} {
	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputTemplate
	k.Config.IncludeFiles = []string{"_loader.html"}

	return ""
}
