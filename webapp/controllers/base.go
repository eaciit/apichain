package controllers

import (
	"github.com/eaciit/orm"
)

type ResultInfo struct {
	IsError bool
	Message string
	Data    interface{}
}

type IBaseController interface {
}

type BaseController struct {
	base IBaseController
	Ctx  *orm.DataContext
}
