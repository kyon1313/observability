package handler

import (
	"net/http"

	"github.com/kyon1313/observability/example/service"
	apw_tracing "github.com/kyon1313/observability/tracing"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service service.UserService
	tracer  apw_tracing.OtelTracing
}

func NewUserHandler(service service.UserService, tracer apw_tracing.OtelTracing) *UserHandler {
	return &UserHandler{service: service, tracer: tracer}
}

func (h *UserHandler) GetUser(ctx *gin.Context) {

	var (
		userName = ctx.Query("name")
		err      error
		response any
		c        = ctx.Request.Context()
	)

	c, span := h.tracer.StartSpan(c, "handler.GetUser")

	response, err = h.service.GetUser(c, userName)
	if err != nil {
		h.tracer.LogTrace(span, &err, "service.GetUser.Error", response)()
		return
	}
	h.tracer.LogTrace(span, &err, "service.Success.Error", response)()
	ctx.JSON(http.StatusOK, response)
}
