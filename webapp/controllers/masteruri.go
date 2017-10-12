package controllers

import (
	"strings"

	db "github.com/eaciit/dbox"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"
	"gopkg.in/mgo.v2/bson"

	. "eaciit/apichain/webapp/helpers"
	. "eaciit/apichain/webapp/models"
)

type MasterUri struct {
	*BaseController
}

func (c *MasterUri) Default(k *knot.WebContext) interface{} {
	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputTemplate
	//k.Config.IncludeFiles = []string{"_loader.html", "resource/ramleditor.html", "resource/resourcemodel.html"}
	return ""
}

func (c *MasterUri) GetData(r *knot.WebContext) interface{} {
	r.Config.OutputType = knot.OutputJson

	crsData, err := c.Ctx.Find(NewUriModel(), tk.M{})
	defer crsData.Close()
	if err != nil {
		return CreateResult(false, nil, err.Error())
	}

	data := []tk.M{}
	err = crsData.Fetch(&data, 0, false)
	if err != nil {
		return CreateResult(false, nil, err.Error())
	}

	return CreateResult(true, data, "")
}

func (c *MasterUri) SaveData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := tk.M{}
	err := k.GetPayload(&p)
	if err != nil {
		return CreateResult(false, nil, err.Error())
	}

	mdl := NewUriModel()
	if p.GetString("_id") != "" {
		mdl.Id = bson.ObjectIdHex(p.GetString("_id"))
	}

	mdl.Name = p.GetString("Name")
	mdl.StringOptions = p.GetString("StringOptions")
	mdl.Options = strings.Split(mdl.StringOptions, "\n")

	err = c.Ctx.Save(mdl)
	if err != nil {
		return CreateResult(false, nil, err.Error())
	}

	return CreateResult(true, nil, "Data has been saved.")
}

func (c *MasterUri) DeleteData(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	p := tk.M{}
	err := k.GetPayload(&p)
	if err != nil {
		return CreateResult(false, nil, err.Error())
	}

	err = c.Ctx.Connection.
		NewQuery().
		From(NewUriModel().TableName()).
		Where(db.Eq("_id", bson.ObjectIdHex(p.GetString("_id")))).
		Delete().
		Exec(nil)
	if err != nil {
		return CreateResult(false, nil, err.Error())
	}

	return CreateResult(true, nil, "Data has been deleted.")
}
