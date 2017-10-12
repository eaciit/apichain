package controllers

import (
	"github.com/eaciit/knot/knot.v1"

	h "eaciit/apichain/webapp/helpers"
)

type Logout struct {
	*BaseController
}

func (c *Logout) Do(k *knot.WebContext) interface{} {
	k.SetSession("userid", nil)
	k.SetSession("username", nil)
	k.SetSession("fullname", nil)
	k.SetSession("usermodel", nil)
	k.SetSession("stime", nil)

	h.Redirect(k, "login", "default")

	return ""
}
