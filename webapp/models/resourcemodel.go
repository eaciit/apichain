package models

import (
	"strings"

	"gopkg.in/mgo.v2/bson"

	"github.com/eaciit/dbox"
	"github.com/eaciit/orm"
	tk "github.com/eaciit/toolkit"

	h "eaciit/apichain/webapp/helpers"
)

type ResourceModel struct {
	orm.ModelBase `bson:"-",json:"-"`

	Id                  bson.ObjectId ` bson:"_id" , json:"_id" `
	Code                string        `bson:"code",json:"Code"`
	Name                string        `bson:"name",json:"Name"`
	Version             string        `bson:"version",json:"Version"`
	Tag                 []string      `bson:"tag",json:"Tag"`
	Description         string        `bson:"description",json:"Description"`
	Raml                int           `bson:"raml",json:"Raml"`
	Stage               string        `bson:"stage",json:"Stage"`
	Uri                 string        `bson:"uri",json:"Uri"`
	Parent              bson.ObjectId ` bson:"parent", json:"parent" `
	ProdCoverageCountry []string      `bson:"prodcoveragecountry",json:"ProdCoverageCountry"`
	ProdCoverageSystem  []string      `bson:"prodcoveragesystem",json:"ProdCoverageSystem"`
	TestCoverageCountry []string      `bson:"testcoveragecountry",json:"TestCoverageCountry"`
	TestCoverageSystem  []string      `bson:"testcoveragesystem",json:"TestCoverageSystem"`
}

func (e *ResourceModel) RecordID() interface{} {
	return e.Id
}

func (m *ResourceModel) TableName() string {
	return "Resource"
}

func (m *ResourceModel) GetAll() ([]tk.M, error) {
	result := []tk.M{}

	c, e := h.PrepareConnection()
	defer c.Close()

	if e != nil {
		return result, e
	}

	csr, e := c.NewQuery().From("Resource").Cursor(nil)
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

func (m *ResourceModel) Delete() error {
	conn, err := h.PrepareConnection()
	defer conn.Close()

	if err != nil {
		return err
	}

	err = conn.NewQuery().
		From(m.TableName()).
		Delete().
		Where(dbox.Eq("parent", m.Parent)).
		Exec(nil)
	if err != nil {
		return err
	}
	err = conn.NewQuery().
		From("ResourceArchive").
		Delete().
		Where(dbox.Eq("parent", m.Parent)).
		Exec(nil)
	if err != nil {
		return err
	}

	return nil
}

func (m *ResourceModel) SaveTo(tableName string) error {
	conn, err := h.PrepareConnection()
	defer conn.Close()
	if err != nil {
		return err
	}

	err = conn.NewQuery().
		From(tableName).
		Save().
		Exec(tk.M{"data": *m})
	if err != nil {
		return err
	}

	return err
}

func (m *ResourceModel) Save() error {
	conn, err := h.PrepareConnection()
	defer conn.Close()
	if err != nil {
		return err
	}

	var isVersionExists = false

	if string(m.Id) == "" {
		m.Id = bson.NewObjectId()
	} else {
		isVersionExists = m.isVersionExist()

		if !isVersionExists {
			//delete exists at resource
			m.DeletebyCode()
			m.Id = bson.NewObjectId()
		}

	}
	//if new version = delete exists version in resource , create new id and insert in both collection

	err = conn.NewQuery().
		From(m.TableName()).
		Save().
		Exec(tk.M{"data": *m})
	if err != nil {
		return err
	}

	err = conn.NewQuery().
		From("ResourceArchive").
		Save().
		Exec(tk.M{"data": *m})
	if err != nil {
		return err
	}

	return nil
}

func (m *ResourceModel) isLastVersion() bool {
	conn, err := h.PrepareConnection()
	defer conn.Close()

	if err != nil {
		return false
	}

	var dbFilter []*dbox.Filter
	dbFilter = append(dbFilter, dbox.Eq("code", m.Code))
	dbFilter = append(dbFilter, dbox.Eq("version", m.Version))

	//filter check where current version and code
	//if false check version is new or not, if new update resource insert archive

	csr, err := conn.NewQuery().
		From(m.TableName()).
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
		return false
	}
	return true
}

func (m *ResourceModel) ListVersionbyCode() ([]string, error) {
	result := []tk.M{}

	var returnVersion = []string{}

	if string(m.Id) == "" {
		return returnVersion, nil
	}

	c, e := h.PrepareConnection()
	defer c.Close()

	if e != nil {
		tk.Println(result, e)
	}
	var dbFilter []*dbox.Filter
	dbFilter = append(dbFilter, dbox.Eq("parent", m.Id))

	csr, e := c.NewQuery().
		From("ResourceArchive").
		Where(dbox.And(dbFilter...)).
		Cursor(nil)
	defer csr.Close()
	if e != nil {
		tk.Println(result, e)
	} else if csr != nil {
		e = csr.Fetch(&result, 0, false)
		if e != nil {
			tk.Println(result, e)
		}

		tk.Println(result, e)
	}

	result2 := []tk.M{}
	csr, e = c.NewQuery().
		From(m.TableName()).
		Where(dbox.And(dbFilter...)).
		Cursor(nil)
	if e != nil {
		tk.Println(result2, e)
	} else if csr != nil {
		e = csr.Fetch(&result2, 0, false)
		if e != nil {
			tk.Println(result2, e)
		}

		tk.Println(result2, e)
	}

	result = append(result, result2...)

	for _, eachVersion := range result {
		varVersion := strings.ToLower(eachVersion.GetString("version"))
		returnVersion = append(returnVersion, varVersion)
	}

	return returnVersion, nil
}

func (m *ResourceModel) CrawlResourcebyCode() (interface{}, error) {
	result := []tk.M{}

	c, e := h.PrepareConnection()
	defer c.Close()

	if e != nil {
		tk.Println(result, e)
	}
	var dbFilter []*dbox.Filter
	dbFilter = append(dbFilter, dbox.Eq("parent", m.Id))

	csr, e := c.NewQuery().
		From("ResourceArchive").
		Where(dbox.And(dbFilter...)).
		Cursor(nil)
	defer csr.Close()
	if e != nil {
		tk.Println(result, e)
	} else if csr != nil {
		e = csr.Fetch(&result, 0, false)
		if e != nil {
			tk.Println(result, e)
		}

		tk.Println(result, e)
	}

	result2 := []tk.M{}
	csr, e = c.NewQuery().
		From(m.TableName()).
		Where(dbox.And(dbFilter...)).
		Cursor(nil)
	if e != nil {
		tk.Println(result2, e)
	} else if csr != nil {
		e = csr.Fetch(&result2, 0, false)
		if e != nil {
			tk.Println(result2, e)
		}

		tk.Println(result2, e)
	}

	result = append(result, result2...)

	return result, nil
}

func (m *ResourceModel) GetResourceArchivebyCodenVersion() ([]tk.M, error) {

	conn, err := h.PrepareConnection()
	defer conn.Close()

	if err != nil {
		tk.Println(nil, err)
	}

	var dbFilter []*dbox.Filter
	dbFilter = append(dbFilter, dbox.Eq("code", m.Code))
	dbFilter = append(dbFilter, dbox.Eq("version", m.Version))

	csr, err := conn.NewQuery().
		From("ResourceArchive").
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

func (m *ResourceModel) isVersionExist() bool {
	conn, err := h.PrepareConnection()
	defer conn.Close()

	if err != nil {
		return false
	}

	var dbFilter []*dbox.Filter
	dbFilter = append(dbFilter, dbox.Eq("code", m.Code))
	dbFilter = append(dbFilter, dbox.Eq("version", m.Version))

	//filter check where current version and code
	//if false check version is new or not, if new update resource insert archive

	csr, err := conn.NewQuery().
		From("ResourceArchive").
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
		return false
	}
	return true
}

func (m *ResourceModel) DeletebyCode() error {
	conn, err := h.PrepareConnection()
	defer conn.Close()

	if err != nil {
		return err
	}

	err = conn.NewQuery().
		From(m.TableName()).
		Delete().
		Where(dbox.Eq("code", m.Code)).
		Exec(nil)
	if err != nil {
		return err
	}

	return nil
}

func (m *ResourceModel) InsertDataPrev() error {
	conn, err := h.PrepareConnection()
	defer conn.Close()

	if err != nil {
		tk.Println(nil, err)
	}

	var dbFilter []*dbox.Filter
	dbFilter = append(dbFilter, dbox.Eq("code", m.Code))

	csr, err := conn.NewQuery().
		From("ResourceArchive").
		Where(dbox.And(dbFilter...)).
		Cursor(nil)
	defer csr.Close()
	if err != nil {
		tk.Println(nil, err)
	}

	data := make([]ResourceModel, 0)
	err = csr.Fetch(&data, 0, false)
	if err != nil {
		tk.Println(nil, err)
	}

	if len(data) == 0 {
		tk.Println(nil, err)
	}
	var lenData = len(data) - 1
	m.Id = data[lenData].Id
	m.Code = data[lenData].Code
	m.Name = data[lenData].Name
	m.Description = data[lenData].Description
	m.Raml = data[lenData].Raml
	m.Version = data[lenData].Version
	m.Stage = data[lenData].Stage
	m.Tag = data[lenData].Tag
	m.Uri = data[lenData].Uri
	m.Save()
	return nil
}
