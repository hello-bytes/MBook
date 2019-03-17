package main

import (
	"fmt"

	"controller/codebook"
	"controller/defaultsite"
	"xweb"

	"app"

	"controller/backend"

	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
)

func startWeb() {
	xweb.InitXWeb(app.RootPath(), app.WebConfigInst().ThemeName(), app.WebConfigInst().CDN())

	e := echo.New()
	e.Renderer = xweb.Templates()

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.BodyLimit("2M"))

	xweb.Instance().SetEcho(e)

	e.Static("/", app.Path("resource/public"))

	defaultsite.NewDefaultController().Init(e, app.WebConfigInst().GetSectionValues("Default"))
	codebook.NewCodeBookController().Init(e, app.WebConfigInst().GetSectionValues("MBook"))
	backend.NewBackendController().Init(e, app.WebConfigInst().GetSectionValues("Backend"))

	port := app.WebConfigInst().Port()
	e.Logger.Fatal(e.Start(fmt.Sprintf(":%d", port)))
}
