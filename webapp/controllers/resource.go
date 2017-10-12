package controllers

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"

	"github.com/eaciit/clit"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/toolkit"
	"gopkg.in/mgo.v2/bson"

	. "eaciit/apichain/webapp/helpers"
	. "eaciit/apichain/webapp/models"
)

type Resource struct {
	*BaseController
}

func (c *Resource) Default(k *knot.WebContext) interface{} {
	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputTemplate
	k.Config.IncludeFiles = []string{"_loader.html", "resource/ramleditor.html", "resource/resourcemodel.html"}
	return ""
}

func (a *Resource) GetAllResource(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	res := toolkit.NewResult()

	varResource := new(ResourceModel)
	varResourceList, err := varResource.GetAll()
	if err != nil {
		res.SetError(err)
		return res
	}

	for i, res := range varResourceList {
		parent := res.Get("parent").(bson.ObjectId)
		id := res.Get("_id").(bson.ObjectId)

		ramlroot := clit.Config("default", "RamlPath", "").(string)
		rootpath := filepath.Join(ramlroot, parent.Hex(), id.Hex())

		count := 0

		_, err = os.Stat(rootpath)

		if !os.IsNotExist(err) {
			filepath.Walk(rootpath, func(path string, f os.FileInfo, err error) error {
				if !f.IsDir() {
					r, err := regexp.MatchString(".raml", f.Name())
					if err == nil && r {
						count++
					}
				}
				return nil
			})
		}

		varResourceList[i].Set("raml", count)
	}

	res.SetData(varResourceList)
	return res
}

func (a *Resource) GetListVersionByCode(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	res := toolkit.NewResult()

	varResource := new(ResourceModel)

	err := k.GetPayload(varResource)
	if err != nil {
		res.SetError(err)
		return res
	}

	listDataVersion, err := varResource.ListVersionbyCode()
	if err != nil {
		res.SetError(err)
		return res
	}

	res.SetData(listDataVersion)
	return res
}

func (a *Resource) GetListResourceByCode(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	res := toolkit.NewResult()

	varResource := new(ResourceModel)

	err := k.GetPayload(varResource)
	if err != nil {
		res.SetError(err)
		return res
	}

	listDataResource, err := varResource.CrawlResourcebyCode()
	if err != nil {
		res.SetError(err)
		return res
	}

	res.SetData(listDataResource)
	return res
}

func (c *Resource) delete(resource *ResourceModel) error {
	err := resource.Delete()
	if err != nil {
		return err
	}

	rootname := clit.Config("default", "RamlPath", "").(string)
	err = os.RemoveAll(filepath.Join(rootname, resource.Parent.Hex()))
	if err != nil {
		return err
	}

	return nil
}

func (c *Resource) Delete(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	res := toolkit.NewResult()

	payload := new(ResourceModel)

	err := k.GetPayload(payload)
	if err != nil {
		res.SetError(err)
		return res
	}

	err = c.delete(payload)
	if err != nil {
		res.SetError(err)
		return res
	}

	return res
}

func (c *Resource) copyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

func (c *Resource) copyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = c.copyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = c.copyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}

func (c *Resource) Save(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	res := toolkit.NewResult()

	payloads := []ResourceModel{}

	err := k.GetPayload(&payloads)
	if err != nil {
		res.SetError(err)
		return res
	}

	oldParent := payloads[0].Parent

	for i, payload := range payloads {
		if string(payload.Id) == "" {
			payloads[i].Id = bson.NewObjectId()
		}
	}

	rootname := clit.Config("default", "RamlPath", "").(string)
	for i, payload := range payloads {
		lastIndex := len(payloads) - 1

		if i == lastIndex {
			payload.Parent = payload.Id
			payloads[i].Parent = payload.Id
			err = payload.SaveTo(payload.TableName())
			if err != nil {
				res.SetError(err)
				return res
			}
		} else {
			payload.Parent = payloads[lastIndex].Id
			payloads[i].Parent = payloads[lastIndex].Id
			err = payload.SaveTo("ResourceArchive")
			if err != nil {
				res.SetError(err)
				return res
			}
		}

		newPath := filepath.Join(rootname, payload.Parent.Hex(), payload.Id.Hex())

		//move resource to new home
		if oldParent.Hex() != "" {
			if oldParent.Hex() != payload.Parent.Hex() {
				errCopy := c.copyDir(filepath.Join(rootname, oldParent.Hex(), payload.Id.Hex()), newPath)
				if errCopy != nil {
					//commonly because new version
					//then copy older version
					c.copyDir(filepath.Join(rootname, oldParent.Hex(), payloads[i-1].Id.Hex()), newPath)
				}
			}
		}

		//make home for homeless resource
		err = os.MkdirAll(newPath, os.ModePerm)
		if err != nil {
			return CreateResult(false, nil, err.Error())
		}
	}

	//delete old home if its not current home
	if oldParent.Hex() != "" {
		if oldParent.Hex() != payloads[0].Parent.Hex() {
			anyPayload := new(ResourceModel)
			anyPayload.Parent = oldParent
			c.delete(anyPayload)
		}
	}

	return CreateResult(true, payloads, "")
}

func (a *Resource) GetDataByCode(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	res := toolkit.NewResult()

	varResource := new(ResourceModel)

	err := k.GetPayload(varResource)
	if err != nil {
		res.SetError(err)
		return res
	}

	listDataVersion, err := varResource.ListVersionbyCode()
	if err != nil {
		res.SetError(err)
		return res
	}

	res.SetData(listDataVersion)
	return res
}
func (a *Resource) GetDataByVersionnCode(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson

	res := toolkit.NewResult()

	varResource := new(ResourceModel)

	err := k.GetPayload(varResource)
	if err != nil {
		res.SetError(err)
		return res
	}

	listDataVersion, err := varResource.GetResourceArchivebyCodenVersion()
	if err != nil {
		res.SetError(err)
		return res
	}

	res.SetData(listDataVersion)
	return res
}
