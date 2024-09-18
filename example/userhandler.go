package main

import (
	"net/http"

	apw_tracing "github.com/kyon1313/observability/tracing"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/codes"
)

type UserHandler struct {
	service UserService
	tracer  apw_tracing.OtelTracing
}

func NewUserHandler(service UserService, tracer apw_tracing.OtelTracing) *UserHandler {
	return &UserHandler{service: service, tracer: tracer}
}

func (h *UserHandler) GetUser(ctx *gin.Context) {
	userName := ctx.Query("name")

	c := ctx.Request.Context()
	c, span := h.tracer.StartSpan(c, "handler.GetUser")
	defer h.tracer.EndSpan(span)

	foundUser, err := h.service.GetUser(c, userName)
	if err != nil {
		h.tracer.SetStatus(span, codes.Error, err.Error())
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, foundUser)
		return
	}

	h.tracer.SetOKStatus(span, "User found")
	ctx.JSON(http.StatusOK, foundUser)
}
