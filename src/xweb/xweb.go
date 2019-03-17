package xweb

import (
	"log"
	"strings"
	"xweb/middleware"

	"xweb/session"

	"github.com/labstack/echo"
)

type XWeb struct {
	cdn       string
	basePath  string
	themeName string
}

var _xweb XWeb

func Instance() *XWeb {
	return &_xweb
}

func InitXWeb(basePath, themeName, cdn string) {
	_xweb.cdn = cdn
	_xweb.basePath = basePath
	_xweb.themeName = themeName
}

func themeHandler(ctx echo.Context) error {
	page := ctx.Request().RequestURI
	page = strings.Replace(page, "//", "/", -1)

	log.Println(page)

	tempArr := strings.Split(page, "/")
	pathArr := make([]string, 0)
	for i := 0; i < len(tempArr); i++ {
		if len(tempArr[i]) > 0 {
			pathArr = append(pathArr, tempArr[i])
		}
	}

	if len(pathArr) <= 2 {
		return RenderStatus(ctx, 404)
	}

	themeName := pathArr[1]
	log.Println(themeName)
	if strings.Compare(themeName, Instance().themeName) != 0 {
		// only request the current theme
		log.Println("theme not match")
		return RenderStatus(ctx, 404)
	}

	path := "/resource/templates/themes"
	for i := 1; i < len(pathArr); i++ {
		if len(pathArr[i]) > 0 {
			path = path + "/" + pathArr[i]
		}
	}

	fullPath := webFile(path)
	index := strings.Index(fullPath, "?")
	if index > 0 {
		fullPath = fullPath[0:index]
	}

	log.Println("the full file is :", fullPath)
	return ctx.File(fullPath)

	//return ctx.String(200, fullPath)
}

func (self *XWeb) SetEcho(e *echo.Echo) {
	session.Init()

	e.GET("/_theme_/*", themeHandler)

	e.HTTPErrorHandler = httpErrorHandler
	e.Use(middleware.AuthCheck("/backend/login"))
	e.Use(middleware.EnvCheck(webFile("resource/config/sk.config"), "/initenv"))
	//e.Use(middleware.EnvCheck(webFile("resource/config/sk")+, "/initenv"))
}

func cdn() string {
	return _xweb.cdn
}

func themeName() string {
	return _xweb.themeName
}

func webFile(relativePath string) string {
	if strings.HasPrefix(relativePath, "/") {
		return _xweb.basePath + relativePath
	} else {
		return _xweb.basePath + "/" + relativePath
	}

}

var debug = true

func httpErrorHandler(err error, c echo.Context) {
	RenderStatus(c, 404)
}
