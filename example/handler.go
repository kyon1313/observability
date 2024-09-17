package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	apw_logging "otel-test/logs"
)

type Handler struct {
	service Service
	logger  apw_logging.OtelLogging
}

func NewHandler(s Service, l apw_logging.OtelLogging) *Handler {
	return &Handler{service: s, logger: l}
}

func (h *Handler) SlowRequest(ctx *gin.Context) {
	c := ctx.Request.Context()
	c, span := h.service.Tracer().Start(c, "slow request")
	defer span.End()

	random, err := h.service.GenerateRandomNumber(c)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error generating random number")
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	sleepDuration := time.Duration(random) * time.Second
	time.Sleep(sleepDuration)

	span.SetAttributes(attribute.String("request_took", fmt.Sprintf("%ds", random)))
	ctx.JSON(http.StatusOK, gin.H{"message": "Slow Request!"})
}

func (h *Handler) PossibleErrorRequest(ctx *gin.Context) {
	c := ctx.Request.Context()
	c, span := h.service.Tracer().Start(c, "possible error request")
	defer span.End()

	random, err := h.service.GenerateRandomNumber(c)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, "error generating random number")
		ctx.JSON(http.StatusInternalServerError, gin.H{"message": "Internal Server Error"})
		return
	}

	if random <= 3 {
		ctx.JSON(http.StatusOK, gin.H{"message": "Good Request!"})
		return
	}

	err = h.service.GenerateError(c)
	span.RecordError(err)
	span.SetStatus(codes.Error, "error encountered")
	span.SetAttributes(attribute.String("error.message", err.Error()))

	ctx.JSON(http.StatusBadRequest, gin.H{"message": "Bad Request!"})
}
