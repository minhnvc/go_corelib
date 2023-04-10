package hserver

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/labstack/echo/v4"
)

type ControllerInterface interface {
	Init(ctx echo.Context)
	Get() error
	Post() error
}

type Controller struct {
	ctx echo.Context
}

// Hàm init đối tượng controller
func (me *Controller) Init(ctx echo.Context) {
	me.ctx = ctx
}

func (me *Controller) GetContext() echo.Context {
	return me.ctx
}

func (me *Controller) Param(key string) string {
	return me.ctx.Param(key)
}

func (me *Controller) SendJSON(obj interface{}) error {
	return me.ctx.JSON(http.StatusOK, obj)
}

func (me *Controller) SendString(obj string) error {
	return me.ctx.String(http.StatusOK, obj)
}

func (me *Controller) GetString(key string) string {
	return me.getData(key)
}

func (me *Controller) GetInt(key string) int {
	result, _ := strconv.Atoi(me.getData(key))
	return result
}

func (me *Controller) GetLong(key string) int64 {
	result, _ := strconv.ParseInt(me.getData(key), 10, 64)
	return result
}

func (me *Controller) GetHeader(key string) string {
	return me.ctx.Request().Header.Get(key)
}

func (me *Controller) SetHeader(key string, value string) {
	me.ctx.Response().Header().Set(key, value)
}

func (me *Controller) GetUserAgent() string {
	return me.GetHeader("user-agent")
}

func (me *Controller) GetIP() string {
	return me.ctx.RealIP()
}

func (me *Controller) GetHostName() string {
	return me.ctx.Request().Host
}

func (me *Controller) BindRequest(obj interface{}) {
	err := me.ctx.Bind(obj)
	if err != nil {
		fmt.Println("Controller", "fail to bind request")
	}
}

func (me *Controller) getData(key string) string {
	if me.ctx.Request().Method == "GET" {
		return me.ctx.QueryParam(key)
	} else if me.ctx.Request().Method == "POST" && me.ctx.Request().Header.Get("content-type") == "application/x-www-form-urlencoded" {
		return me.ctx.FormValue(key)
	}
	return ""
}
