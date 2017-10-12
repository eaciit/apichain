package controllers

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/eaciit/clit"
	"github.com/eaciit/knot/knot.v1"
	tk "github.com/eaciit/toolkit"

	. "eaciit/apichain/webapp/helpers"
)

type Designer struct {
	*BaseController
	rootname string
}

type FileStructure struct {
	Path     string          `json:"path"`
	Name     string          `json:"name"`
	Type     string          `json:"type"`
	Children []FileStructure `json:"children"`
}

type dir struct {
	path string
	os.FileInfo
}

const TYPE_FOLDER = "folder"
const TYPE_FILE = "file"

func (c *Designer) list(parentFS *FileStructure, parentPath string) error {
	files, err := ioutil.ReadDir(parentPath)
	if err != nil {
		return err
	}

	fileInfos := []os.FileInfo{}
	for _, file := range files {
		fileInfos = append(fileInfos, file)
	}

	calcPath := func(path string) string {
		rootPath := filepath.Join("..", c.rootname)
		path = strings.Replace(path, rootPath, "", -1)

		if path == "" && parentPath == rootPath {
			path = "/"
		}

		return path
	}

	calcType := func(f os.FileInfo) string {
		if f.IsDir() {
			return TYPE_FOLDER
		}
		return TYPE_FILE
	}

	for _, fileInfo := range fileInfos {
		if fileInfo.Name() == ".DS_Store" {
			continue
		}

		f := FileStructure{}
		f.Path = calcPath(filepath.Join(parentPath, fileInfo.Name()))
		f.Name = fileInfo.Name()
		f.Type = calcType(fileInfo)
		f.Children = []FileStructure{}

		if fileInfo.IsDir() {
			c.list(&f, filepath.Join(parentPath, fileInfo.Name()))
		}

		parentFS.Children = append(parentFS.Children, f)
	}

	return nil
}

func (c *Designer) List(k *knot.WebContext) interface{} {
	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputJson

	t := tk.M{}
	err := k.GetPayload(&t)
	if err != nil {
		return CreateResult(false, nil, err.Error())
	}

	rootDir := FileStructure{}
	rootDir.Path = "/"
	rootDir.Name = ""
	rootDir.Type = TYPE_FOLDER
	rootDir.Children = []FileStructure{}

	c.rootname = "webapp/ramlfiles" + t.GetString("path")
	rootpath := filepath.Join("..", c.rootname)

	err = c.list(&rootDir, rootpath)
	if err != nil {
		return CreateResult(false, nil, err.Error())
	}

	return CreateResult(true, rootDir, "")
}

func (c *Designer) Save(k *knot.WebContext) interface{} {
	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputJson

	t := tk.M{}
	err := k.GetPayload(&t)
	if err != nil {
		return CreateResult(false, nil, err.Error())
	}
	_ = os.Remove(filepath.Join("..", c.rootname, t.GetString("path")))

	path := filepath.Join("..", c.rootname, t.GetString("path"))
	tk.Println(path)

	_, err = os.Stat(path)

	if os.IsNotExist(err) {
		var file, err = os.Create(path)
		if err != nil {
			return CreateResult(false, nil, err.Error())
		}
		defer file.Close()
	}

	file, err := os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		return CreateResult(false, nil, err.Error())
	}
	defer file.Close()

	// write some text line-by-line to file
	_, err = file.WriteString(t.GetString("contents"))
	if err != nil {
		return CreateResult(false, nil, err.Error())
	}

	// save changes
	err = file.Sync()
	if err != nil {
		return CreateResult(false, nil, err.Error())
	}

	return CreateResult(true, t, "")
}

func (c *Designer) newFolder(path string) error {
	return os.MkdirAll(path, os.ModePerm)
}

func (c *Designer) NewFolder(k *knot.WebContext) interface{} {
	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputJson

	t := tk.M{}
	err := k.GetPayload(&t)
	if err != nil {
		return CreateResult(false, nil, err.Error())
	}

	err = c.newFolder(filepath.Join("..", c.rootname, t.GetString("path")))
	if err != nil {
		return CreateResult(false, nil, err.Error())
	}

	return CreateResult(true, t, "")
}

func (c *Designer) Rename(k *knot.WebContext) interface{} {
	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputJson

	t := tk.M{}
	err := k.GetPayload(&t)
	if err != nil {
		return CreateResult(false, nil, err.Error())
	}

	source := filepath.Join("..", c.rootname, t.GetString("source"))
	destination := filepath.Join("..", c.rootname, t.GetString("destination"))

	err = os.Rename(source, destination)
	if err != nil {
		return CreateResult(false, nil, err.Error())
	}

	return CreateResult(true, t, "")
}

func (c *Designer) Delete(k *knot.WebContext) interface{} {
	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputJson

	t := tk.M{}
	err := k.GetPayload(&t)
	if err != nil {
		return CreateResult(false, nil, err.Error())
	}

	err = os.Remove(filepath.Join("..", c.rootname, t.GetString("path")))
	if err != nil {
		return CreateResult(false, nil, err.Error())
	}

	return CreateResult(true, t, "")
}

func (c *Designer) EditLine(k *knot.WebContext) interface{} {
	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputJson

	t := tk.M{}
	err := k.GetPayload(&t)
	if err != nil {
		return CreateResult(false, nil, err.Error())
	}

	input, err := ioutil.ReadFile(filepath.Join("..", c.rootname, t.GetString("path")))
	if err != nil {
		return CreateResult(false, nil, err.Error())
	}

	lines := strings.Split(string(input), "\n")

	for i, _ := range lines {
		if i == t.GetInt("line")-1 {
			lines[i] = t.GetString("to")
		}
	}

	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(filepath.Join("..", c.rootname, t.GetString("path")), []byte(output), 0644)
	if err != nil {
		return CreateResult(false, nil, err.Error())
	} else {
		return CreateResult(true, nil, "")
	}

	// return CreateResult(false, nil, "Line not found")
}

func (c *Designer) GetRamlList(k *knot.WebContext) interface{} {
	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputJson

	t := tk.M{}
	err := k.GetPayload(&t)
	if err != nil {
		return CreateResult(false, nil, err.Error())
	}

	ramlroot := clit.Config("default", "RamlPath", "").(string)
	rootpath := filepath.Join(ramlroot, t.GetString("parent"), t.GetString("id"))

	paths := []string{}

	err = filepath.Walk(rootpath, func(path string, f os.FileInfo, err error) error {
		if !f.IsDir() {
			r, err := regexp.MatchString(".raml", f.Name())
			if err == nil && r {
				p := strings.Replace(path, ramlroot, "", 1)
				paths = append(paths, p)
			} else {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return CreateResult(false, nil, err.Error())
	}

	return CreateResult(true, paths, "")
}

func (c *Designer) NewRaml(k *knot.WebContext) interface{} {
	k.Config.NoLog = true
	k.Config.OutputType = knot.OutputJson

	t := tk.M{}
	err := k.GetPayload(&t)
	if err != nil {
		return CreateResult(false, nil, err.Error())
	}

	ramlroot := clit.Config("default", "RamlPath", "").(string)
	rootpath := filepath.Join(ramlroot, t.GetString("parent"), t.GetString("id"))
	f, err := os.Create(filepath.Join(rootpath, t.GetString("filename")+".raml"))

	defer f.Close()
	if err != nil {
		return CreateResult(false, nil, err.Error())
	}

	_, err = f.WriteString("#%RAML 1.0\ntitle: \"" + t.GetString("title") + "\"\n\n")
	if err != nil {
		return CreateResult(false, nil, err.Error())
	}
	f.Sync()

	_, err = f.WriteString("securitySchemes:\n")
	if err != nil {
		return CreateResult(false, nil, err.Error())
	}
	f.Sync()

	securitySchemes := t.Get("securitySchemes").([]interface{})
	for _, schemes := range securitySchemes {
		s, _ := tk.ToM(schemes)
		_, err = f.WriteString("  " + s.GetString("name") + ":\n    type:\n")
		if err != nil {
			return CreateResult(false, nil, err.Error())
		}
		f.Sync()
	}

	_, err = f.WriteString("  " + t.GetString("security") + "\n")
	if err != nil {
		return CreateResult(false, nil, err.Error())
	}
	f.Sync()

	_, err = f.WriteString("traceability:\n")
	if err != nil {
		return CreateResult(false, nil, err.Error())
	}
	f.Sync()

	_, err = f.WriteString("  " + t.GetString("traceability") + "\n\n")
	if err != nil {
		return CreateResult(false, nil, err.Error())
	}
	f.Sync()

	selectedSchema := t.Get("selectedSchema").([]interface{})

	if len(selectedSchema) != 0 {
		_, err = f.WriteString("types:\n")
		if err != nil {
			return CreateResult(false, nil, err.Error())
		}
		f.Sync()

		err = c.newFolder(filepath.Join(rootpath, "schema"))
		if err != nil {
			return CreateResult(false, nil, err.Error())
		}

		for _, schema := range selectedSchema {
			sch, _ := tk.ToM(schema)

			//create schema file
			fschema, err := os.Create(filepath.Join(rootpath, "schema", sch.GetString("name")+".json"))
			defer fschema.Close()
			if err != nil {
				return CreateResult(false, nil, err.Error())
			}

			jsonbyte, _ := json.MarshalIndent(sch.Get("body"), "", "  ")

			_, err = fschema.WriteString(string(jsonbyte))
			if err != nil {
				return CreateResult(false, nil, err.Error())
			}
			fschema.Sync()

			_, err = f.WriteString("  " + sch.GetString("name") + ": " + "!include schema/" + sch.GetString("name") + ".json\n")
			if err != nil {
				return CreateResult(false, nil, err.Error())
			}
			f.Sync()
		}

		_, err = f.WriteString("\n")
		if err != nil {
			return CreateResult(false, nil, err.Error())
		}
		f.Sync()
	}

	_, err = f.WriteString("/" + t.GetString("resourceName") + ":\n")
	if err != nil {
		return CreateResult(false, nil, err.Error())
	}
	f.Sync()

	_, err = f.WriteString("  description: \"" + t.GetString("description") + "\"\n")
	if err != nil {
		return CreateResult(false, nil, err.Error())
	}
	f.Sync()

	isGet := t.Get("isMethodGet").(bool)

	if isGet {
		_, err = f.WriteString("  get:\n    responses:\n")
		if err != nil {
			return CreateResult(false, nil, err.Error())
		}
		f.Sync()

		responses := t.Get("methodGetResponses").([]interface{})
		for _, response := range responses {
			r, _ := tk.ToM(response)
			_, err = f.WriteString("      " + r.GetString("code") + ":\n        body:\n          application/json:\n            schema: " + r.GetString("body") + "\n")
			if err != nil {
				return CreateResult(false, nil, err.Error())
			}
			f.Sync()
		}

	}

	isPost := t.Get("isMethodPost").(bool)

	if isPost {
		_, err = f.WriteString("  post:\n    responses:\n")
		if err != nil {
			return CreateResult(false, nil, err.Error())
		}
		f.Sync()

		responses := t.Get("methodPostResponses").([]interface{})
		for _, response := range responses {
			r, _ := tk.ToM(response)
			_, err = f.WriteString("      " + r.GetString("code") + ":\n        body:\n          application/json:\n            schema: " + r.GetString("body") + "\n")
			if err != nil {
				return CreateResult(false, nil, err.Error())
			}
			f.Sync()
		}

	}

	return CreateResult(true, nil, "")
}
