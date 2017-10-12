package models

import (
	//"strings"

	"gopkg.in/mgo.v2/bson"

	"github.com/eaciit/dbox"
	"github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"

	h "eaciit/apichain/webapp/helpers"
)

type SchemaModel struct {
	orm.ModelBase `bson:"-",json:"-"`

	Id         bson.ObjectId ` bson:"_id" , json:"_id" `
	Name       string        `bson:"name",json:"name"`
	Jsonstring string        `bson:"jsonstring",json:"jsonstring"`
}

func (e *SchemaModel) RecordID() interface{} {
	return e.Id
}

func (m *SchemaModel) TableName() string {
	return "Schema"
}

func (m *SchemaModel) GetAll() ([]tk.M, error) {
	result := []tk.M{}

	c, e := h.PrepareConnection()
	defer c.Close()

	if e != nil {
		return result, e
	}

	csr, e := c.NewQuery().From("Schema").Cursor(nil)
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

func (m *SchemaModel) Delete() error {
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

func (m *SchemaModel) Save() error {
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
func (m *SchemaModel) GetDatabyName() ([]tk.M, error) {

	conn, err := h.PrepareConnection()
	defer conn.Close()

	if err != nil {
		tk.Println(nil, err)
	}

	var dbFilter []*dbox.Filter
	dbFilter = append(dbFilter, dbox.Eq("code", m.Name))
	//dbFilter = append(dbFilter, dbox.Eq("version", m.Version))

	csr, err := conn.NewQuery().
		From("Schema").
		Where(dbox.And(dbFilter...)).
		Cursor(nil)
	defer csr.Close()
	if err != nil {
		tk.Println(nil, err)
	}

	data := make([]tk.M, 0)
	err = csr.Fetch(&data, 0, false)
	if err != nil {
		tk.Println(nil, err)
	}

	if len(data) == 0 {
		tk.Println(nil, err)
	}
	return data, err
}
