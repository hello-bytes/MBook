package defaultsite

import "github.com/labstack/echo"

type DefaultController struct {
}

var defaultController *DefaultController

func NewDefaultController() *DefaultController {
	if defaultController == nil {
		defaultController = &DefaultController{}
	}
	return defaultController
}

func (self *DefaultController) Init(e *echo.Echo, params map[string]string) {
	registerRoute(e)
}
