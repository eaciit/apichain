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
	"io"
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

	ramlroot := clit.Config("default", "RamlPath", "").(string)
	rootpath := filepath.Join(ramlroot, t.GetString("parent"), t.GetString("id"))

	var list func(parentFS *FileStructure, parentPath string) error
	list = func(parentFS *FileStructure, parentPath string) error {
		var files []os.FileInfo
		files, err = ioutil.ReadDir(parentPath)
		if err != nil {
			return err
		}

		fileInfos := []os.FileInfo{}
		for _, file := range files {
			fileInfos = append(fileInfos, file)
		}

		calcPath := func(path string) string {
			path = strings.Replace(path, rootpath, "", -1)

			if path == "" && parentPath == rootpath {
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
				list(&f, filepath.Join(parentPath, fileInfo.Name()))
			}

			parentFS.Children = append(parentFS.Children, f)
		}

		return nil
	}

	err = list(&rootDir, rootpath)
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

	ramlroot := clit.Config("default", "RamlPath", "").(string)
	rootpath := filepath.Join(ramlroot, t.GetString("parent"), t.GetString("id"))
	fpath := filepath.Join(rootpath, t.GetString("path"))

	os.Remove(fpath)

	_, err = os.Stat(fpath)

	if os.IsNotExist(err) {
		var file *os.File
		file, err = os.Create(fpath)
		defer file.Close()
		if err != nil {
			return CreateResult(false, nil, err.Error())
		}
	}

	file, err := os.OpenFile(fpath, os.O_RDWR, 0644)
	defer file.Close()
	if err != nil {
		return CreateResult(false, nil, err.Error())
	}

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

	ramlroot := clit.Config("default", "RamlPath", "").(string)
	rootpath := filepath.Join(ramlroot, t.GetString("parent"), t.GetString("id"))

	err = c.newFolder(filepath.Join(rootpath, t.GetString("path")))
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

	ramlroot := clit.Config("default", "RamlPath", "").(string)
	rootpath := filepath.Join(ramlroot, t.GetString("parent"), t.GetString("id"))

	source := filepath.Join(rootpath, t.GetString("source"))
	destination := filepath.Join(rootpath, t.GetString("destination"))

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

	ramlroot := clit.Config("default", "RamlPath", "").(string)
	rootpath := filepath.Join(ramlroot, t.GetString("parent"), t.GetString("id"))
	fpath := filepath.Join(rootpath, t.GetString("path"))

	err = os.Remove(fpath)
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

func (c *Designer) writeJSONFile(path string, jsonbody map[string]interface{}) error {
	//create schema file
	fschema, err := os.Create(path)
	defer fschema.Close()
	if err != nil {
		return err
	}

	jsonbyte, _ := json.MarshalIndent(jsonbody, "", "  ")
	_, err = fschema.WriteString(string(jsonbyte))
	if err != nil {
		return err
	}
	fschema.Sync()

	return nil
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
	newramlpath := filepath.Join(rootpath, t.GetString("filename")+".raml")

	f, err := os.Create(newramlpath)

	defer f.Close()
	if err != nil {
		return CreateResult(false, nil, err.Error())
	}

	toWrite := ""

	toWrite += "#%RAML 1.0\n"
	toWrite += "title: \"" + t.GetString("filename") + "\"\n"
	toWrite += "baseUri: " + t.GetString("baseUri") + "\n\n"

	toWrite += "securitySchemes:\n"

	securitySchemes := t.Get("securitySchemes").([]interface{})
	for _, schemes := range securitySchemes {
		s, _ := tk.ToM(schemes)

		if s.GetString("name") != "" {
			toWrite += "  " + s.GetString("name") + ":\n"
			toWrite += "    description: |\n      " + s.GetString("description") + "\n"
			toWrite += "    type: " + s.GetString("type") + "\n"

			if s.Get("body") != nil {
				toWrite += "    describedBy:\n"
				toWrite += "      headers:\n"
				toWrite += "        Authorization:\n"
				toWrite += "          description: |\n"
				toWrite += "            Used to send a valid OAuth 2 access token. Do not use\n"
				toWrite += "            with the 'access_token' query string parameter.\n"
				toWrite += "          type: string\n"
				toWrite += "      queryParameters:\n"
				toWrite += "        access_token:\n"
				toWrite += "          description: |\n"
				toWrite += "            Used to send a valid OAuth 2 access token. Do not use together with\n"
				toWrite += "            the 'Authorization' header\n"
				toWrite += "          type: string\n"
				toWrite += "      responses:\n"

				toBody := s.Get("body").(map[string]interface{})
				for i, toBodyDetail := range toBody {
					if i == "responses" {
						for _, response := range toBodyDetail.([]interface{}) {
							var detailResponse = response.(map[string]interface{})
							toWrite += "        " + detailResponse["code"].(string) + ":\n"
							toWrite += "          description: |\n"
							toWrite += "            " + detailResponse["description"].(string) + "\n"

						}
					}
				}

				var authorizationUri = toBody["authorizationUri"].(string)
				var accessTokenUri = toBody["accessTokenUri"].(string)
				var authorizationGrants = toBody["authorizationGrants"].(string)

				toWrite += "    settings:\n"
				toWrite += "      authorizationUri: " + authorizationUri + "\n"
				toWrite += "      accessTokenUri: " + accessTokenUri + "\n"
				toWrite += "      authorizationGrants: " + authorizationGrants + "\n"
			}
		}
	}
	toWrite += "\n"

	selectedSchema := t.Get("selectedSchema").([]interface{})
	if len(selectedSchema) != 0 {
		toWrite += "types:\n"

		err = c.newFolder(filepath.Join(rootpath, "schema"))
		if err != nil {
			return CreateResult(false, nil, err.Error())
		}

		for _, schema := range selectedSchema {
			sch, _ := tk.ToM(schema)

			var jsonbody map[string]interface{}
			json.Unmarshal([]byte(sch.GetString("body")), &jsonbody)

			err = c.writeJSONFile(filepath.Join(rootpath, "schema", sch.GetString("name")+".json"), jsonbody)
			if err != nil {
				return CreateResult(false, nil, err.Error())
			}

			toWrite += "  " + sch.GetString("name") + ": " + "!include schema/" + sch.GetString("name") + ".json\n"
		}

		toWrite += "\n"
	}

	if t.GetString("resourceName") != "" {
		toWrite += "/" + t.GetString("resourceName") + ":\n"
		toWrite += "  description: \"" + t.GetString("description") + "\"\n"

		isGet := t.Get("isMethodGet").(bool)
		if isGet {
			toWrite += "  get:\n"

			toWrite += "    headers:\n"
			traceabilities := t.Get("methodGetTraceabilities").([]interface{})
			for _, traceability := range traceabilities {
				r, _ := tk.ToM(traceability)

				toWrite += "      " + r.GetString("name") + ": " + r.GetString("body") + "\n"
			}
			toWrite += "    responses:\n"

			responses := t.Get("methodGetResponses").([]interface{})
			for _, response := range responses {
				r, _ := tk.ToM(response)

				toWrite += "      " + r.GetString("code") + ":\n        body:\n          application/json:\n            schema: " + r.GetString("body") + "\n"
			}

		}

		isPost := t.Get("isMethodPost").(bool)
		if isPost {
			toWrite += "  post:\n"

			toWrite += "    headers:\n"
			traceabilities := t.Get("methodPostTraceabilities").([]interface{})
			for _, traceability := range traceabilities {
				r, _ := tk.ToM(traceability)

				toWrite += "      " + r.GetString("name") + ": " + r.GetString("body") + "\n"
			}
			toWrite += "    responses:\n"

			responses := t.Get("methodPostResponses").([]interface{})
			for _, response := range responses {
				r, _ := tk.ToM(response)

				toWrite += "      " + r.GetString("code") + ":\n        body:\n          application/json:\n            schema: " + r.GetString("body") + "\n"
			}

		}
	}

	_, err = f.WriteString(toWrite)
	if err != nil {
		return CreateResult(false, nil, err.Error())
	}
	f.Sync()

	return CreateResult(true, newramlpath, "")
}

func (c *Designer) Upload(k *knot.WebContext) interface{} {
	k.Config.OutputType = knot.OutputJson
	ramlPath := clit.Config("default", "RamlPath", "").(string)
	parentPath := k.Request.FormValue("parent")
	idPath := k.Request.FormValue("id")
	additionalPath := k.Request.FormValue("additionalPath")
	file, handler, err := k.Request.FormFile("batchFile")
	defer file.Close()
	if err != nil {

	}
	uploadfilename := handler.Filename

	dstSource := ramlPath + tk.PathSeparator + parentPath + tk.PathSeparator + idPath + tk.PathSeparator + additionalPath + tk.PathSeparator + uploadfilename

	f, err := os.OpenFile(dstSource, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {

	}
	defer f.Close()
	io.Copy(f, file)
	if err != nil {
		//c.ErrorResultInfo("Err", err.Error())
	}
	return nil
}
