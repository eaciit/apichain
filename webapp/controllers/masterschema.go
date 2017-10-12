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

type MasterSchema struct {
	*BaseController
}

func (c *MasterSchema) Default(k *knot.WebContext) interface{} {
	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputTemplate
	//k.Config.IncludeFiles = []string{"_loader.html", "resource/ramleditor.html", "resource/resourcemodel.html"}
	return ""
}

func (a *MasterSchema) GetAllSchema(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	res := toolkit.NewResult()

	varSchema := new(SchemaModel)
	varSchemaList, err := varSchema.GetAll()
	if err != nil {
		res.SetError(err)
		return res
	}

	res.SetData(varSchemaList)
	return res
}

func (c *MasterSchema) Delete(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	res := toolkit.NewResult()

	payload := new(SchemaModel)

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

func (c *MasterSchema) Save(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	res := toolkit.NewResult()

	payload := new(SchemaModel)

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
func (a *MasterSchema) GetDatabyName(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	res := toolkit.NewResult()

	varSchema := new(SchemaModel)
	varSchemaList, err := varSchema.GetDatabyName()
	if err != nil {
		res.SetError(err)
		return res
	}

	res.SetData(varSchemaList)
	return res
}
