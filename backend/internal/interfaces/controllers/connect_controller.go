package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"

	"babywolf-blog-backend/internal/application/connect"
)

// ConnectController expone el canal de MobileConnect (Práctica 6).
type ConnectController struct {
	store *connect.Store
}

func NewConnectController(store *connect.Store) *ConnectController {
	return &ConnectController{store: store}
}

type enviarConnectRequest struct {
	Emisor      string `json:"emisor"`
	Dispositivo string `json:"dispositivo"`
	Texto       string `json:"texto"`
}

func (c *ConnectController) Enviar(ctx *gin.Context) {
	var req enviarConnectRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Cuerpo de la petición inválido."})
		return
	}
	m, err := c.store.Enviar(ctx.Param("sala"), req.Emisor, req.Dispositivo, req.Texto)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusCreated, m)
}

func (c *ConnectController) Listar(ctx *gin.Context) {
	desde, _ := strconv.ParseInt(ctx.Query("desde"), 10, 64)
	msgs, conteo, err := c.store.Listar(ctx.Param("sala"), desde)
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"mensajes": msgs, "conteo": conteo})
}

func (c *ConnectController) Limpiar(ctx *gin.Context) {
	if err := c.store.Limpiar(ctx.Param("sala")); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}
	ctx.Status(http.StatusNoContent)
}
