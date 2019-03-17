package codebook

import (
	"html/template"
	"io/ioutil"
	"log"
	"net/url"
	"strings"

	"xweb"

	gomd "github.com/hello-bytes/GoMD"
	"github.com/labstack/echo"
)

func isImage(filePath string) bool {
	if strings.HasSuffix(filePath, ".png") || strings.HasSuffix(filePath, ".jpg") ||
		strings.HasSuffix(filePath, ".jpeg") || strings.HasSuffix(filePath, ".gif") || strings.HasSuffix(filePath, ".bmp") {
		return true
	}
	return false
}

func codebook(ctx echo.Context) error {
	page := urlToPath(ctx.Request().RequestURI)
	page = strings.Replace(page, "//", "/", -1)
	tempArr := strings.Split(page, "/")

	pathArr := make([]string, 0)
	for i := 0; i < len(tempArr); i++ {
		if len(tempArr[i]) > 0 {
			pathArr = append(pathArr, tempArr[i])
		}
	}

	if len(pathArr) < 1 {
		return xweb.RenderStatus(ctx, 404)
	}

	realtivePath := ""
	bookName := pathArr[0]
	for i := 1; i < len(pathArr); i++ {
		realtivePath = realtivePath + "/" + pathArr[i]
	}

	log.Println("get the path", realtivePath)
	return bookProcess(ctx, bookName, realtivePath)
}

func smartReadContent(fullPath string) ([]byte, error) {
	contentbin, err := ioutil.ReadFile(fullPath)
	if err != nil {
		if strings.HasSuffix(fullPath, ".html") {
			fullPath = fullPath[0 : len(fullPath)-4]
			fullPath = fullPath + "md"
			contentbin, err = ioutil.ReadFile(fullPath)
			if err != nil {
				return nil, err
			}
		}
	}
	return contentbin, err
}

func urlToPath(requestUrl string) string {
	uri, _ := url.Parse(requestUrl)
	return uri.Path
}

func bookProcess(c echo.Context, bookName, relativePath string) error {
	if len(relativePath) == 0 {
		relativePath = "README.md"
	}

	bookPath := codeBookController.BookPath(bookName)
	if bookPath == nil {
		log.Println("Not Found==" + bookName + ":" + relativePath)
		return xweb.RenderStatus(c, 404)
	}

	page := urlToPath(c.Request().RequestURI)

	fullPath := *bookPath + "/" + relativePath
	log.Println("the full path :", fullPath)

	// 如果是图片,直接返回文件
	if isImage(fullPath) {
		return c.File(fullPath)
	}

	pageData := make(map[string]interface{})
	pageData["CatalogueHtml"] = template.HTML(codeBookController.catalogueMgr().renderRootNode(bookName, page))

	contentbin, err := smartReadContent(fullPath) //ioutil.ReadFile(fullPath)

	if err == nil {
		htmlContent := gomd.Run(contentbin)
		pageData["ArticleContent"] = template.HTML(htmlContent)
	} else {
		pageData["ArticleContent"] = "Can not read this document : " + fullPath
	}
	return xweb.RenderArticle(c, "pages/article.html", pageData)
}

func registerRoute(e *echo.Echo) {
	bookList := codeBookController.BookItemList()
	for _, v := range bookList {
		bookObj := v
		log.Println("register the book request :", bookObj.id)
		e.GET("/"+bookObj.id+"*", codebook)
	}
}
