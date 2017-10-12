package models

import (
	"github.com/eaciit/orm"
	"gopkg.in/mgo.v2/bson"
)

type UriModel struct {
	orm.ModelBase `bson:"-",json:"-"`
	Id            bson.ObjectId `bson:"_id" , json:"_id" `
	Name          string        `bson:"Name",json:"Name"`
	Options       []string      `bson:"Options",json:"Options"`
	StringOptions string        `bson:"StringOptions",json:"StringOptions"`
}

func NewUriModel() *UriModel {
	m := new(UriModel)
	m.Id = bson.NewObjectId()
	return m
}

func (e *UriModel) RecordID() interface{} {
	return e.Id
}

func (m *UriModel) TableName() string {
	return "MasterUri"
}
