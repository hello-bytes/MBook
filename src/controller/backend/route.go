package backend

import (
	"app"
	"log"
	"os/exec"
	"strings"
	"xweb"

	"io/ioutil"

	"xweb/session"

	"github.com/labstack/echo"
)

func initEnvPage(ctx echo.Context) error {
	app.MakeSecurityKey()
	pageData := make(map[string]interface{}, 0)
	pageData["Sk"] = app.SecurityKey()
	pageData["SkFilePath"] = app.SkFile()
	session.GetSession(ctx.Request()).Set("__role__", "|admin|init|")
	return xweb.RenderPage(ctx, "backend/initenv.html", pageData)
}

func loginPage(ctx echo.Context) error {
	log.Println("login page....")

	// if it is areadly loged, redirect to /backend/index
	role := session.GetSession(ctx.Request()).GetString("__role__")
	log.Println("role = ", role)
	if strings.Contains(role, "admin") {
		return ctx.Redirect(301, "/backend/index")
	}

	return xweb.RenderBlankPage(ctx, "backend/login.html", nil)
}
func doLoginPage(ctx echo.Context) error {
	sk := ctx.FormValue("sk")
	if strings.Compare(sk, app.SecurityKey()) == 0 {
		session.GetSession(ctx.Request()).Set("__role__", "|admin|init|")

		return ctx.Redirect(301, "/backend/index")
	}
	return ctx.Redirect(301, "/backend/login")
}

func logoutPage(ctx echo.Context) error {
	session.GetSession(ctx.Request()).Delete("__role__")
	return ctx.Redirect(301, "/backend/login")
}

func indexPage(ctx echo.Context) error {
	return xweb.RenderBackend(ctx, "backend/index.html", nil)
}

func operatorPage(ctx echo.Context) error {
	return xweb.RenderBackend(ctx, "backend/operator.html", nil)
}

func restartPage(ctx echo.Context) error {
	return xweb.RenderBackend(ctx, "backend/restart.html", nil)
}

func doRestartPage(ctx echo.Context) error {
	cmd := "cd " + app.RootPath() + " && bash deploy.sh restart"
	exec_shell_no_result(cmd)
	return ctx.String(200, "ok=")
}

func checkAlive(ctx echo.Context) error {
	return ctx.String(200, "alive")
}

func doUpdateConfig(ctx echo.Context) error {
	configText := ctx.FormValue("config")
	configFilePath := app.Path("resource/config/env.config")
	ioutil.WriteFile(configFilePath, []byte(configText), 0x755)
	return ctx.Redirect(301, "/backend/config")
}

func configPage(ctx echo.Context) error {
	configText := ""

	filePath := app.Path("resource/config/env.config")
	skbin, err := ioutil.ReadFile(filePath)
	if err == nil {
		configText = string(skbin)
	}

	data := make(map[string]interface{}, 0)
	data["ConfigText"] = configText
	return xweb.RenderBackend(ctx, "backend/config.html", data)
}

func exec_shell_no_result(command string) {
	cmd := exec.Command("/bin/bash", "-c", command)
	cmd.Run()
}

func registerRoute(e *echo.Echo) {
	e.GET("/initenv", initEnvPage)

	e.GET("/backend/login", loginPage)
	e.POST("/backend/login.do", doLoginPage)
	e.GET("/backend/logoutnow", logoutPage)

	e.GET("/backend/index", indexPage)
	e.GET("/backend/operator", operatorPage)

	e.GET("/backend/restart", restartPage)
	e.GET("/backend/restart.do", doRestartPage)

	e.GET("/backend/config", configPage)
	e.POST("/backend/updateconfig.do", doUpdateConfig)

	e.GET("/api/checkalive", checkAlive)
}
