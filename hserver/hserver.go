package hserver

import (
	"reflect"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type HServer struct {
	echo *echo.Echo
}

func New() *HServer {
	hserver := &HServer{}
	hserver.Init()
	return hserver
}

func (me *HServer) Init() {
	me.echo = echo.New()
	me.echo.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"http://*", "https://*", "zbrowser://*"},
	}))
}

// Thêm controller vào route của webserver
func (me *HServer) Route(path string, handler ControllerInterface) {
	me.echo.Any(path, func(ctx echo.Context) error {
		reflectVal := reflect.ValueOf(handler)
		t := reflect.Indirect(reflectVal).Type()
		vc := reflect.New(t)
		exec := vc.Interface().(ControllerInterface)
		exec.Init(ctx)
		if ctx.Request().Method == "GET" {
			return exec.Get()
		} else if ctx.Request().Method == "POST" {
			return exec.Post()
		}
		return nil
	})
}

func (me *HServer) RouteStatic(path string, localPath string) {
	me.echo.Static(path, localPath)
}

func (me *HServer) RouteStaticFile(path string, filePath string) {
	me.echo.File(path, filePath)
}

func (me *HServer) RouteStaticFileNoCache(path string, filePath string) {
	me.echo.File(path, filePath, noCache)
}

func (me *HServer) StartServer(port int) {
	me.echo.Logger.Fatal(me.echo.Start(":" + strconv.Itoa(port)))
}

func noCache(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Add("Cache-Control", "no-cache, private, max-age=0")
		c.Response().Header().Add("Pragma", "no-cache")
		c.Response().Header().Add("X-Accel-Expires", "0")
		next(c)
		return nil
	}
}
