package xweb

import (
	"bytes"
	"html/template"
	"io"
	"log"
	"net/http"
	"path/filepath"
	"strconv"

	"github.com/labstack/echo"
)

//============================================================
// ******** Templates ********
//============================================================

type PluginCallHandler interface {
	Call(interface{}) map[string]interface{}
}

var pluginCalls map[string]PluginCallHandler

func RenderBackend(c echo.Context, contentTemplate string, data map[string]interface{}) error {
	return renderImpl(c, http.StatusOK, "backend", contentTemplate, data)
}

func RenderArticle(c echo.Context, contentTemplate string, data map[string]interface{}) error {
	return renderImpl(c, http.StatusOK, "article", contentTemplate, data)
}

func RenderPage(c echo.Context, contentTemplate string, data map[string]interface{}) error {
	return renderImpl(c, http.StatusOK, "index", contentTemplate, data)
}

func RenderBlankPage(c echo.Context, contentTemplate string, data map[string]interface{}) error {
	return renderImpl(c, http.StatusOK, "blank", contentTemplate, data)
}

func RenderStatus(c echo.Context, httpStatus int) error {
	page := strconv.Itoa(httpStatus)
	return renderImpl(c, httpStatus, "index", "status/"+page+".html", nil)
}

func renderImpl(c echo.Context, httpStatus int, layoutTemplate, contentTemplate string, data map[string]interface{}) error {
	if data == nil {
		data = make(map[string]interface{}, 0)
	}
	data["pageTemplate"] = contentTemplate
	data["__RequestUrl__"] = c.Request().RequestURI
	return c.Render(httpStatus, layoutTemplate, data)
}

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
	return t.templates.ExecuteTemplate(w, name, data)
}

func tplInclude(file string, dot map[string]interface{}) template.HTML {
	var buffer = &bytes.Buffer{}
	templateFile := webFile("resource/templates/themes/" + themeName() + "/" + file)
	//log.Println("include file", templateFile)
	tpl, err := template.New(filepath.Base(file)).Funcs(funcMap).ParseFiles(templateFile)
	if err != nil {
		log.Println("get tpl error for file(", file, ") error:", err)
		return "" //"Can not found file : " + file
	}
	err = tpl.Execute(buffer, dot)
	if err != nil {
		log.Println("Execute template file(", file, ") syntax error:", err)
		return ""
	}

	return template.HTML(buffer.String())
}

func RegisterPluginCall(name string, handler PluginCallHandler) {
	if pluginCalls == nil {
		pluginCalls = make(map[string]PluginCallHandler, 0)
	}
	if len(name) == 0 {
		return
	}
	pluginCalls[name] = handler
}

func pluginCall(pluginId string, params interface{}) map[string]interface{} {
	handler, ok := pluginCalls[pluginId]
	if ok {
		return handler.Call(params)
	}

	result := make(map[string]interface{}, 0)
	result["code"] = 1
	return result
}

var funcMap = template.FuncMap{
	"localAssets": func(path string) string {
		return path + "?v=1027"
	},
	"cdnAssets": func(path string) string {
		return cdn() + path + "?v=1024"
	},
	"themeAssets": func(path string) string {
		return "/_theme_/" + themeName() + "/" + path + "?v=1027"
	},
	"pluginCall": func(pluginId string, params interface{}) map[string]interface{} {
		return pluginCall(pluginId, params)
	},
}

func Templates() *Template {
	//log.Println("theme name", themeName())
	templates := template.New(webFile("resource/templates/themes/" + themeName() + "/layout/*.html")).Funcs(funcMap).Funcs(template.FuncMap{"include": tplInclude})
	templates.ParseGlob(webFile("resource/templates/themes/" + themeName() + "/layout/backend/*.html"))
	templates.ParseGlob(webFile("resource/templates/themes/" + themeName() + "/layout/frontend/*.html"))
	allTemplates := &Template{
		templates: templates,
	}

	log.Println(webFile("resource/templates/themes/" + themeName() + "/layout/frontend/*.html"))

	return allTemplates
}
