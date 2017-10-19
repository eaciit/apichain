package controllers

import (
	"time"

	db "github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"

	h "eaciit/apichain/webapp/helpers"
	m "eaciit/apichain/webapp/models"
)

type Login struct {
	*BaseController
}

func (c *Login) Default(k *knot.WebContext) interface{} {
	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputTemplate
	k.Config.LayoutTemplate = ""
	k.Config.IncludeFiles = []string{}
	if k.Session("userid") == nil {
		k.Session("userid", nil)
	}

	return ""
}

func (c *Login) Do(k *knot.WebContext) interface{} {
	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputJson

	formdata := struct {
		Username string
		Password string
	}{}

	err := k.GetPayload(&formdata)
	if err != nil {
		return h.ErrorResult(err)
	}

	conn, err := h.PrepareConnection()
	defer conn.Close()
	if err != nil {
		return h.CreateResult(false, nil, err.Error())
	}

	filter := []*db.Filter{}
	filter = append(filter, db.Eq("username", formdata.Username))

	csr, err := conn.NewQuery().From("SysUsers").Where(filter...).Cursor(nil)
	defer csr.Close()
	if err != nil {
		return h.CreateResult(false, nil, err.Error())
	}

	res := make([]m.SysUserModel, 0)

	err = csr.Fetch(&res, 0, false)
	if err != nil {
		return h.ErrorResult(err)
	}

	if len(res) > 0 {
		resUser := res[0]
		if h.GetSHA256(formdata.Password) == resUser.Password {
			if resUser.Enable == true {
				k.SetSession("userid", resUser.Id.Hex())
				k.SetSession("username", resUser.Username)
				k.SetSession("fullname", resUser.Fullname)
				k.SetSession("usermodel", resUser)
				k.SetSession("stime", time.Now())

				return h.CreateResult(true, nil, "Login successful.")
			} else {
				return h.CreateResult(false, nil, "Your account is disabled, please contact administrator to enable it.")
			}
		} else {
			return h.CreateResult(false, nil, "Invalid Username or password!")
		}
	} else {
		return h.CreateResult(false, nil, "Invalid Username or password!")
	}
}
