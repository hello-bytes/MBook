package defaultsite

import (
	"controller/codebook"

	"xweb"

	"github.com/labstack/echo"
)

func aboutPage(c echo.Context) error {
	pageData := make(map[string]interface{})
	return xweb.RenderPage(c, "pages/about.html", pageData)
}

func index(c echo.Context) error {
	pageData := make(map[string]interface{})
	pageData["CodeBookList"] = codebook.NewCodeBookController().GetCodeBookList()

	return xweb.RenderPage(c, "pages/home.html", pageData)
}

func registerRoute(e *echo.Echo) {
	e.GET("/", index)
	e.GET("/about.html", aboutPage)
}
