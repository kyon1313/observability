package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service UserService
}

func NewUserHandler(service UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) AddUser(ctx *gin.Context) {
	var newUser User

	c := ctx.Request.Context()
	c, span := h.service.Tracer().Start(c, "handler.AddUser")
	defer span.End()

	if err := ctx.ShouldBindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := h.service.AddUser(c, newUser); err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add user"})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{"message": "User added successfully"})
}
func (h *UserHandler) GetUser(ctx *gin.Context) {
	userName := ctx.Query("name")

	c := ctx.Request.Context()
	c, span := h.service.Tracer().Start(c, "handler.GetUser")
	defer span.End()

	foundUser, err := h.service.GetUser(c, userName)
	if err != nil {
		ctx.Error(err)
		ctx.JSON(http.StatusBadRequest, foundUser)
		return
	}

	ctx.JSON(http.StatusOK, foundUser)
}
