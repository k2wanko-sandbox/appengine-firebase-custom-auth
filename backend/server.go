package backend

import (
	"crypto/rsa"
	"encoding/json"
	"net/http"
	"os"
	"text/template"
	"time"

	"golang.org/x/net/context"

	"github.com/dgrijalva/jwt-go"
	"github.com/k2wanko/echo-appengine"
	"github.com/labstack/echo"
	"github.com/labstack/echo/engine/standard"
	"github.com/mjibson/goon"
	"google.golang.org/appengine/log"
)

var (
	serviceAccountEmail                 = ""
	privateKey          *rsa.PrivateKey = nil
)

func init() {
	// time
	loc, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	time.Local = loc

	// load service account
	svf, err := os.Open("service.json")
	if err != nil {
		panic(err)
	}

	sv := &struct {
		Email      string `json:"client_email"`
		PrivateKey string `json:"private_key"`
	}{}
	if err := json.NewDecoder(svf).Decode(sv); err != nil {
		panic(err)
	}
	serviceAccountEmail = sv.Email
	if key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(sv.PrivateKey)); err != nil {
		panic(err)
	} else {
		privateKey = key
	}

	// setup framework
	e := echo.New()
	e.Use(appengine.AppContext())

	// set reder
	t := &Template{
		templates: template.Must(template.ParseGlob("public/*.html")),
	}
	e.SetRenderer(t)

	// handle
	e.GET("/", handle)
	e.POST("/login", handleLogin)

	// set handle
	s := standard.New("")
	s.SetHandler(e)
	http.Handle("/", s)
}

func handle(c echo.Context) error {
	return c.Render(200, "index", "")
}

func handleLogin(c echo.Context) error {
	email := c.FormValue("email")
	logf(c, "handleCreate: email=%s", email)

	if u, _ := userFromEmail(c, email); u != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "dup email")
	}

	u := &User{
		Email: email,
	}
	g := goon.FromContext(c)
	_, err := g.Put(u)
	if err != nil {
		return err
	}

	tok := createToken(u)
	s, err := tok.SignedString(privateKey)
	if err != nil {
		return err
	}

	reft, err := u.RefreshToken()
	if err != nil {
		return err
	}

	r := struct {
		User         *User
		Token        string
		RefreshToken string
	}{
		User:         u,
		Token:        s,
		RefreshToken: reft,
	}

	return c.JSON(http.StatusCreated, r)
}

func logf(ctx context.Context, formant string, args ...interface{}) {
	log.Infof(ctx, formant, args...)
}
