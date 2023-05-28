package routes

import (
	"github.com/github.com/vido21/dating-app/auth"
	"github.com/github.com/vido21/dating-app/blogs"
	"github.com/github.com/vido21/dating-app/common"
	premiumPackages "github.com/github.com/vido21/dating-app/premium-packages"
	"github.com/github.com/vido21/dating-app/profiles"
	"github.com/github.com/vido21/dating-app/purchases"
	"github.com/labstack/echo/v4"
)

func DefineApiRoute(e *echo.Echo) {
	controllers := []common.Controller{
		auth.AuthController{},
		blogs.BlogsController{},
		profiles.ProfilesController{},
		premiumPackages.PremiumPackageController{},
		purchases.PurchaseController{},
	}
	var routes []common.Route
	for _, controller := range controllers {
		routes = append(routes, controller.Routes()...)
	}
	api := e.Group("/api/v0")
	for _, route := range routes {
		switch route.Method {
		case echo.POST:
			{
				api.POST(route.Path, route.Handler, route.Middleware...)
				break
			}
		case echo.GET:
			{
				api.GET(route.Path, route.Handler, route.Middleware...)
				break
			}
		case echo.DELETE:
			{
				api.DELETE(route.Path, route.Handler, route.Middleware...)
				break
			}
		case echo.PUT:
			{
				api.PUT(route.Path, route.Handler, route.Middleware...)
				break
			}
		case echo.PATCH:
			{
				api.PATCH(route.Path, route.Handler, route.Middleware...)
				break
			}
		}
	}
}
