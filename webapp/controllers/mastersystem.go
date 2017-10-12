package controllers

import (
	//"os"
	//"path/filepath"

	//"github.com/eaciit/clit"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
	//"gopkg.in/mgo.v2/bson"

	//. "eaciit/apichain/webapp/helpers"
	. "eaciit/apichain/webapp/models"
)

type MasterSystem struct {
	*BaseController
}

func (c *MasterSystem) Default(k *knot.WebContext) interface{} {
	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputTemplate
	//k.Config.IncludeFiles = []string{"_loader.html", "resource/ramleditor.html", "resource/resourcemodel.html"}
	return ""
}

func (a *MasterSystem) GetAllSystem(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	res := toolkit.NewResult()

	varSystem := new(SystemModel)
	varSystemList, err := varSystem.GetAll()
	if err != nil {
		res.SetError(err)
		return res
	}

	res.SetData(varSystemList)
	return res
}

func (c *MasterSystem) Delete(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	res := toolkit.NewResult()

	payload := new(SystemModel)

	err := k.GetPayload(payload)
	if err != nil {
		res.SetError(err)
		return res
	}

	err = payload.Delete()
	if err != nil {
		res.SetError(err)
		return res
	}

	return res
}

func (c *MasterSystem) Save(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	res := toolkit.NewResult()

	payload := new(SystemModel)

	err := k.GetPayload(payload)
	if err != nil {
		res.SetError(err)
		return res
	}

	err = payload.Save()
	if err != nil {
		res.SetError(err)
		return res
	}

	return res
}
