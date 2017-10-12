package models

import (
	//"strings"

	"gopkg.in/mgo.v2/bson"

	"github.com/eaciit/dbox"
	"github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"

	h "eaciit/apichain/webapp/helpers"
)

type HttpStatusesModel struct {
	orm.ModelBase `bson:"-",json:"-"`

	Id   bson.ObjectId ` bson:"_id" , json:"_id" `
	Name string        `bson:"name",json:"name"`
	Code string        `bson:"code",json:"code"`
}

func (e *HttpStatusesModel) RecordID() interface{} {
	return e.Id
}

func (m *HttpStatusesModel) TableName() string {
	return "HttpStatuses"
}

func (m *HttpStatusesModel) GetAll() ([]tk.M, error) {
	result := []tk.M{}

	c, e := h.PrepareConnection()
	defer c.Close()

	if e != nil {
		return result, e
	}

	csr, e := c.NewQuery().From("HttpStatuses").Cursor(nil)
	defer csr.Close()
	if e != nil {
		return result, e
	} else if csr != nil {
		defer csr.Close()

		e = csr.Fetch(&result, 0, false)

		if e != nil {

			return result, e
		}

		return result, nil
	}

	return result, nil
}

func (m *HttpStatusesModel) Delete() error {
	conn, err := h.PrepareConnection()
	defer conn.Close()

	if err != nil {
		return err
	}

	err = conn.NewQuery().
		From(m.TableName()).
		Delete().
		Where(dbox.Eq("_id", m.Id)).
		Exec(nil)
	if err != nil {
		return err
	}

	return nil
}

func (m *HttpStatusesModel) Save() error {
	conn, err := h.PrepareConnection()
	defer conn.Close()
	if err != nil {
		return err
	}

	if string(m.Id) == "" {
		m.Id = bson.NewObjectId()
	}
	//if new version = delete exists version in resource , create new id and insert in both collection

	err = conn.NewQuery().
		From(m.TableName()).
		Save().
		Exec(tk.M{"data": *m})
	if err != nil {
		return err
	}

	return nil
}