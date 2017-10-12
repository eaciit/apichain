package helper

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"net/http"

	"github.com/eaciit/clit"
	db "github.com/eaciit/dbox"
	_ "github.com/eaciit/dbox/dbc/mongo"
	knot "github.com/eaciit/knot/knot.v1"
)

var (
	DebugMode bool
)

func PrepareConnection() (db.IConnection, error) {
	ci := &db.ConnectionInfo{
		Host:     clit.Config("default", "Host", "").(string),
		Database: clit.Config("default", "Database", "").(string),
		UserName: clit.Config("default", "Username", "").(string),
		Password: clit.Config("default", "Password", "").(string),
		Settings: nil,
	}
	fmt.Printf("%+v\n", ci)

	c, e := db.NewConnection("mongo", ci)

	if e != nil {
		return nil, e
	}

	e = c.Connect()
	if e != nil {
		return nil, e
	}

	return c, nil
}

func IsAuthenticate(k *knot.WebContext) {
	if k.Session("userid") == nil {
		Redirect(k, "login", "default")
	}
	return
}

func Redirect(k *knot.WebContext, controller string, action string) {
	http.Redirect(k.Writer, k.Request, "/"+controller+"/"+action, http.StatusTemporaryRedirect)
}

func CreateResult(success bool, data interface{}, message string) map[string]interface{} {
	if !success {
		fmt.Println("ERROR! ", message)
		if DebugMode {
			panic(message)
		}
	}

	return map[string]interface{}{
		"data":    data,
		"success": success,
		"message": message,
	}
}

func ErrorResult(err error) map[string]interface{} {
	return CreateResult(false, nil, err.Error())
}

func GetSHA256(str string) string {
	sha256Bytes := sha256.Sum256([]byte(str))
	return hex.EncodeToString(sha256Bytes[:])
}
