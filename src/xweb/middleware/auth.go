package middleware

import (
	"log"
	"net/http"
	"strings"
	"sync"

	"xweb/session"

	"app"

	"github.com/labstack/echo"
)

var urlAuthMap sync.Map

func initAuthMap() {
	urlAuthMap.Store("/initenv", "init")
	urlAuthMap.Store("/backend/index", "admin")
	urlAuthMap.Store("/backend/restart", "admin")
	urlAuthMap.Store("/backend/restart.do", "admin")
	urlAuthMap.Store("/backend/config", "admin")
	urlAuthMap.Store("/backend/operator", "admin")

	urlAuthMap.Store("/backend/updateconfig.do", "admin")

}

func urlToAuth(url string) string {
	index := strings.Index(url, "?")
	if index > 0 {
		url = url[0:index]
	}
	result, ok := urlAuthMap.Load(url)
	if ok {
		return stringVal(result, "")
	}
	return "none"
}

func stringVal(str interface{}, def string) string {
	switch str.(type) {
	case string:
		return str.(string)
	default:
		return def
	}
}

type AuthChecker struct {
	forbiddenUrl string
}

func AuthCheck(forbiddenUrl string) echo.MiddlewareFunc {
	initAuthMap()
	auth := AuthChecker{forbiddenUrl: forbiddenUrl}
	return authMiddleware(auth)
}

func authMiddleware(authConfig AuthChecker) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set(echo.HeaderServer, "MBook Server")

			url := c.Request().RequestURI
			authKey := urlToAuth(url)
			session := session.SessionStart(c.Response().Writer, c.Request())

			if app.WebConfigInst().IsDebug() {
				return next(c)
			}

			//如果找不到对Auth Key，则直接返回true
			if len(authKey) == 0 || strings.Compare("none", authKey) == 0 {
				return next(c)
			}

			log.Println("check auth key [" + authKey + "] for [" + url + "]")

			if strings.Compare(authKey, "init") == 0 {
				// 初始化页面，并且没有被初始化，则跳转
				if !app.HasSecurityKey() {
					return next(c)
				}
			}

			checkKey := "|" + authKey + "|"
			authorizeKey := session.GetString("__role__")
			if strings.Contains(authorizeKey, checkKey) || strings.Compare(authorizeKey, authKey) == 0 {
				log.Println("pass, the auth key is " + authorizeKey)
				return next(c)
			}

			log.Println(url)

			//redirecturl := "/backend/login"
			redirecturl := authConfig.forbiddenUrl
			log.Println("check failed, will redirect to [", redirecturl, "]")
			return c.Redirect(http.StatusTemporaryRedirect, redirecturl)
		}
	}
}
