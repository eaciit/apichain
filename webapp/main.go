package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/eaciit/clit"
	_ "github.com/eaciit/dbox/dbc/mongo"
	"github.com/eaciit/knot/knot.v1"
	"github.com/eaciit/orm"

	c "eaciit/apichain/webapp/controllers"
	h "eaciit/apichain/webapp/helpers"
)

var (
	err error
)

func main() {
	clit.LoadConfigFromFlag("", "", filepath.Join(clit.ExeDir(), "..", "config", "app.json"))
	if err = clit.Commit(); err != nil {
		kill(err)
	}
	defer clit.Close()

	app := createApp()
	webhost := clit.Config("default", "WebHost", "").(string)

	routes := map[string]knot.FnContent{
		"/": func(r *knot.WebContext) interface{} {
			h.Redirect(r, "login", "default")
			return true
		},
		"prerequest": func(r *knot.WebContext) interface{} {
			url := r.Request.URL.String()

			if url == "/login/default" && r.Session("username") != nil {
				h.Redirect(r, "dashboard", "default")
				return true
			}

			if strings.Index(url, "/login") < 0 && url != "/" {
				h.IsAuthenticate(r)
				return nil
			}
			return nil
		},
		"postrequest": func(r *knot.WebContext) interface{} {
			return nil
		},
	}

	knot.StartAppWithFn(app, webhost, routes)
}

func createApp() *knot.App {
	webroot := clit.Config("default", "WebRoot", "").(string)
	if webroot == "" {
		webroot = clit.ExeDir()
	}

	conn, err := h.PrepareConnection()
	if err != nil {
		kill(err)
	}

	baseCtrl := new(c.BaseController)
	ctx := orm.New(conn)
	baseCtrl.Ctx = ctx

	app := knot.NewApp("apichain")

	/**REGISTER ALL CONTROLLERS HERE**/
	app.Register(&c.Dashboard{BaseController: baseCtrl})
	app.Register(&c.Designer{BaseController: baseCtrl})
	app.Register(&c.Login{BaseController: baseCtrl})
	app.Register(&c.Logout{BaseController: baseCtrl})
	app.Register(&c.Resource{BaseController: baseCtrl})
	app.Register(&c.MasterStage{BaseController: baseCtrl})
	app.Register(&c.MasterCountry{BaseController: baseCtrl})
	app.Register(&c.MasterSystem{BaseController: baseCtrl})
	app.Register(&c.MasterUri{BaseController: baseCtrl})
	app.Register(&c.MasterSchema{BaseController: baseCtrl})
	app.Register(&c.MasterHttpStatuses{BaseController: baseCtrl})

	/* FOLDER STRUCTURE */
	app.Static("libs", filepath.Join(webroot, "assets", "core", "lib"))
	app.Static("styles", filepath.Join(webroot, "assets", "core", "apps", "css"))
	app.Static("scripts", filepath.Join(webroot, "assets", "core", "apps", "js"))
	app.Static("images", filepath.Join(webroot, "assets", "img"))
	app.Static("plugins", filepath.Join(webroot, "assets", "vendors"))
	app.Static("modules", filepath.Join(webroot, "assets", "node_modules"))
	app.Static("files", filepath.Join(webroot, "ramlfiles"))
	/* END FOLDER STRUCTURE */

	app.LayoutTemplate = "_layout.html"
	app.ViewsPath = filepath.Join(webroot, "views")

	return app
}

func kill(err error) {
	fmt.Printf("error. %s \n", err.Error())
	os.Exit(100)
}
