package backend

import (
	"github.com/labstack/echo"
)

type BackendController struct {
}

var backendController *BackendController

func NewBackendController() *BackendController {
	if backendController == nil {
		backendController = &BackendController{}
	}
	return backendController
}

func (self *BackendController) Init(e *echo.Echo, params map[string]string) {
	registerRoute(e)
}
