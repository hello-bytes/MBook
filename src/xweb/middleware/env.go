package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/labstack/echo"
)

type EnvChecker struct {
	_file string
	_url  string
}

var _webRoot string

func fileExist(filename string) bool {
	_, err := os.Stat(filename)
	return err == nil || os.IsExist(err)
}

func EnvCheck(file, url string) echo.MiddlewareFunc {
	envCheck := EnvChecker{_file: file, _url: url}
	return envCheckWithConfig(envCheck)
}

func envCheckWithConfig(envCheck EnvChecker) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			url := c.Request().RequestURI
			if strings.Compare(url, envCheck._url) == 0 {
				return next(c)
			}

			if fileExist(envCheck._file) {
				return next(c)
			} else {
				return c.Redirect(http.StatusTemporaryRedirect, envCheck._url)
			}
		}
	}
}
