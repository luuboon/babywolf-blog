package controllers

import (
	"crypto/rand"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
	"unicode/utf8"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"babywolf-blog-backend/internal/application/soporte"
	"babywolf-blog-backend/internal/application/usecases"
	"babywolf-blog-backend/internal/domain/entities"
	"babywolf-blog-backend/internal/infrastructure/database"
	"babywolf-blog-backend/internal/infrastructure/repositories"
	"babywolf-blog-backend/internal/infrastructure/security"
)

// Mesa de ayuda BabyWolf — Práctica 7 de Desarrollo móvil.
//
// Roles: admin (entra con su cuenta del blog), agente y usuarios (entran con
// una URL preasignada que genera el admin). Todas las rutas /api/soporte/*
// usan el encabezado X-Soporte-Token que se entrega al unirse.

const (
	propositoInvitacion = "soporte-invitacion"
	ttlInvitacion       = 30 * time.Minute
	ttlSesionSoporte    = 12 * time.Hour
	bucketSoporte       = "soporte_archivos"
	maxDispositivo      = 40
)

type SoporteController struct {
	repo  *repositories.GormSoporteRepository
	users *usecases.ManageUsersUseCase
}

func NewSoporteController(repo *repositories.GormSoporteRepository, users *usecases.ManageUsersUseCase) *SoporteController {
	return &SoporteController{repo: repo, users: users}
}

// ── utilidades ─────────────────────────────────────────────────

func nuevoUUID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	b[6] = (b[6] & 0x0f) | 0x40
	b[8] = (b[8] & 0x3f) | 0x80
	h := hex.EncodeToString(b)
	return h[0:8] + "-" + h[8:12] + "-" + h[12:16] + "-" + h[16:20] + "-" + h[20:32]
}

func nuevoID(n int) string {
	b := make([]byte, n)
	_, _ = rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

func limpiarDispositivo(d string) string {
	d = strings.TrimSpace(d)
	if d == "" {
		return "Teléfono"
	}
	if utf8.RuneCountInString(d) > maxDispositivo {
		d = string([]rune(d)[:maxDispositivo])
	}
	return d
}

// responder traduce errores a respuestas HTTP. Los errores internos se
// registran en el log del servidor y al cliente sólo le llega un mensaje
// genérico, sin detalles técnicos.
func responder(ctx *gin.Context, err error) {
	var v *soporte.ErrValidacion
	switch {
	case errors.As(err, &v):
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "Revisa los datos marcados.", "campos": v.Campos})
	case errors.Is(err, soporte.ErrSinPermiso):
		ctx.JSON(http.StatusForbidden, gin.H{"error": "No tienes permiso para esta acción."})
	case errors.Is(err, gorm.ErrRecordNotFound):
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No se encontró el registro."})
	default:
		log.Printf("[soporte] %s %s: %v", ctx.Request.Method, ctx.FullPath(), err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "No fue posible completar la operación."})
	}
}

type estadoParticipante struct {
	entities.SoporteParticipante
	Conectado   bool   `json:"conectado"`
	RolEfectivo string `json:"rol_efectivo"`
}

func (c *SoporteController) estado(ps []entities.SoporteParticipante, ahora time.Time) []estadoParticipante {
	out := make([]estadoParticipante, 0, len(ps))
	for _, p := range ps {
		out = append(out, estadoParticipante{p, soporte.Conectado(p, ahora), soporte.RolEfectivo(p, ps, ahora)})
	}
	return out
}

func yo(ctx *gin.Context) entities.SoporteParticipante {
	return ctx.MustGet("soporte_yo").(entities.SoporteParticipante)
}

func rolEf(ctx *gin.Context) string {
	return ctx.GetString("soporte_rol")
}

// Auth valida X-Soporte-Token, registra el latido y calcula el rol efectivo.
func (c *SoporteController) Auth() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		id, err := security.LeerToken(ctx.GetHeader("X-Soporte-Token"), time.Now())
		if err != nil {
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Tu sesión no es válida. Vuelve a unirte."})
			return
		}
		ahora := time.Now().UTC()
		if err := c.repo.Latido(ctx.Request.Context(), id, &ahora); err != nil {
			responder(ctx, err)
			ctx.Abort()
			return
		}
		ps, err := c.repo.Participantes(ctx.Request.Context())
		if err != nil {
			responder(ctx, err)
			ctx.Abort()
			return
		}
		var actual *entities.SoporteParticipante
		for i := range ps {
			if ps[i].ID == id {
				actual = &ps[i]
			}
		}
		if actual == nil { // expulsado (activo = false) o inexistente
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Tu sesión no es válida. Vuelve a unirte."})
			return
		}
		ctx.Set("soporte_yo", *actual)
		ctx.Set("soporte_todos", ps)
		ctx.Set("soporte_rol", soporte.RolEfectivo(*actual, ps, ahora))
		ctx.Next()
	}
}

// ── sesión ─────────────────────────────────────────────────────

type dispositivoRequest struct {
	Dispositivo string `json:"dispositivo"`
}

// SesionAdmin: el admin entra con su JWT de Supabase (cuenta del blog). El
// rol se comprueba contra la base de datos, no contra lo que diga el token.
func (c *SoporteController) SesionAdmin(ctx *gin.Context) {
	uid, _ := ctx.Get("user_id")
	uidStr, _ := uid.(string)
	user, err := c.users.GetUserBySupabaseUID(ctx.Request.Context(), uidStr)
	if err != nil || user == nil || !user.IsAdmin() || !user.Active {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "Sólo un administrador del blog puede abrir la mesa de ayuda."})
		return
	}
	var req dispositivoRequest
	_ = ctx.ShouldBindJSON(&req)
	ahora := time.Now().UTC()

	p, err := c.repo.AdminDeUsuario(ctx.Request.Context(), user.ID)
	switch {
	case errors.Is(err, gorm.ErrRecordNotFound):
		nombre := user.Username
		if utf8.RuneCountInString(nombre) > soporte.MaxNombre {
			nombre = string([]rune(nombre)[:soporte.MaxNombre])
		}
		uidDB := user.ID
		p = &entities.SoporteParticipante{
			ID: nuevoUUID(), Nombre: nombre, Rol: soporte.RolAdmin, UserID: &uidDB,
			Dispositivo: limpiarDispositivo(req.Dispositivo), Activo: true, UltimoLatido: &ahora,
		}
		if err := c.repo.CrearParticipante(ctx.Request.Context(), p); err != nil {
			responder(ctx, err)
			return
		}
	case err != nil:
		responder(ctx, err)
		return
	default:
		p.Dispositivo = limpiarDispositivo(req.Dispositivo)
		p.Activo = true
		p.UltimoLatido = &ahora
		if err := c.repo.ActualizarParticipante(ctx.Request.Context(), p.ID, map[string]any{
			"dispositivo": p.Dispositivo, "activo": true, "ultimo_latido": ahora,
		}); err != nil {
			responder(ctx, err)
			return
		}
	}
	c.entregarSesion(ctx, *p, http.StatusOK)
}

func (c *SoporteController) entregarSesion(ctx *gin.Context, p entities.SoporteParticipante, status int) {
	token, err := security.EmitirToken(p.ID, ttlSesionSoporte, time.Now())
	if err != nil {
		responder(ctx, err)
		return
	}
	ctx.JSON(status, gin.H{"token": token, "participante": p})
}

// Invitar genera la URL preasignada con la que un dispositivo se une con un
// rol fijo. La URL vale 30 minutos y sólo se puede usar una vez.
func (c *SoporteController) Invitar(ctx *gin.Context) {
	if !soporte.PuedeAdministrar(rolEf(ctx)) {
		responder(ctx, soporte.ErrSinPermiso)
		return
	}
	var req struct {
		Rol    string `json:"rol"`
		Nombre string `json:"nombre"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Cuerpo de la petición inválido."})
		return
	}
	if err := soporte.ValidarRolInvitacion(req.Rol); err != nil {
		responder(ctx, err)
		return
	}
	nombre, err := soporte.ValidarNombre(req.Nombre)
	if err != nil {
		responder(ctx, err)
		return
	}
	ahora := time.Now()
	invID := nuevoID(9)
	q, expira, err := security.Firmar(propositoInvitacion, url.Values{
		"inv": {invID}, "rol": {req.Rol}, "nombre": {nombre},
	}, ttlInvitacion, ahora)
	if err != nil {
		responder(ctx, err)
		return
	}
	inv := &entities.SoporteInvitacion{
		ID: invID, Rol: req.Rol, Nombre: nombre, ExpiraEn: expira, CreadaPor: yo(ctx).ID,
	}
	if err := c.repo.CrearInvitacion(ctx.Request.Context(), inv); err != nil {
		responder(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, gin.H{
		"url":    baseURL(ctx) + "/api/soporte/unirse?" + q.Encode(),
		"rol":    req.Rol,
		"nombre": nombre,
		"expira": expira.Format(time.RFC3339),
	})
}

// verificarInvitacion valida la firma de la URL y devuelve la invitación.
func (c *SoporteController) verificarInvitacion(ctx *gin.Context) (*entities.SoporteInvitacion, bool) {
	q := ctx.Request.URL.Query()
	if err := security.Verificar(propositoInvitacion, q, time.Now()); err != nil {
		if errors.Is(err, security.ErrExpirada) {
			ctx.JSON(http.StatusGone, gin.H{"error": "La invitación expiró. Pide una nueva al administrador."})
		} else {
			ctx.JSON(http.StatusForbidden, gin.H{"error": "La URL de invitación no es válida o fue modificada."})
		}
		return nil, false
	}
	inv, err := c.repo.Invitacion(ctx.Request.Context(), q.Get("inv"))
	if err != nil {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "La URL de invitación no es válida o fue modificada."})
		return nil, false
	}
	return inv, true
}

// InfoInvitacion (GET) deja ver qué rol otorga la URL sin consumirla.
func (c *SoporteController) InfoInvitacion(ctx *gin.Context) {
	inv, ok := c.verificarInvitacion(ctx)
	if !ok {
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"rol": inv.Rol, "nombre": inv.Nombre, "expira": inv.ExpiraEn.Format(time.RFC3339), "usada": inv.UsadaPor != nil,
	})
}

// Unirse (POST a la misma URL preasignada) crea el participante y entrega su
// token. Consumir la invitación es atómico (ver repositorio).
func (c *SoporteController) Unirse(ctx *gin.Context) {
	inv, ok := c.verificarInvitacion(ctx)
	if !ok {
		return
	}
	if inv.UsadaPor != nil {
		ctx.JSON(http.StatusConflict, gin.H{"error": "Esta invitación ya fue usada en otro dispositivo."})
		return
	}
	var req dispositivoRequest
	_ = ctx.ShouldBindJSON(&req)
	ahora := time.Now().UTC()
	p := &entities.SoporteParticipante{
		ID: nuevoUUID(), Nombre: inv.Nombre, Rol: inv.Rol,
		Dispositivo: limpiarDispositivo(req.Dispositivo), Activo: true, UltimoLatido: &ahora,
	}
	if err := c.repo.Unirse(ctx.Request.Context(), inv.ID, p, ahora); err != nil {
		if errors.Is(err, repositories.ErrInvitacionUsada) {
			ctx.JSON(http.StatusConflict, gin.H{"error": "Esta invitación ya fue usada en otro dispositivo."})
			return
		}
		responder(ctx, err)
		return
	}
	c.entregarSesion(ctx, *p, http.StatusCreated)
}

// Latido mantiene la conexión y devuelve el estado de los 4 dispositivos.
func (c *SoporteController) Latido(ctx *gin.Context) {
	todos := ctx.MustGet("soporte_todos").([]entities.SoporteParticipante)
	ctx.JSON(http.StatusOK, gin.H{
		"yo":            yo(ctx),
		"rol_efectivo":  rolEf(ctx),
		"participantes": c.estado(todos, time.Now().UTC()),
	})
}

// Salir es la desconexión explícita: el dispositivo deja de contar como
// conectado y, si era el admin, el agente pasa a cubrirlo.
func (c *SoporteController) Salir(ctx *gin.Context) {
	if err := c.repo.Latido(ctx.Request.Context(), yo(ctx).ID, nil); err != nil {
		responder(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// Expulsar da de baja a un participante (sólo admin real).
func (c *SoporteController) Expulsar(ctx *gin.Context) {
	if rolEf(ctx) != soporte.RolAdmin {
		responder(ctx, soporte.ErrSinPermiso)
		return
	}
	if ctx.Param("id") == yo(ctx).ID {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": "No puedes expulsarte a ti mismo."})
		return
	}
	if err := c.repo.ActualizarParticipante(ctx.Request.Context(), ctx.Param("id"), map[string]any{"activo": false}); err != nil {
		responder(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// Monitor es la comprobación en tiempo real: sólo admin y agente. La app la
// consulta cada pocos segundos y recarga cuando algo cambió.
func (c *SoporteController) Monitor(ctx *gin.Context) {
	if !soporte.PuedeMonitorear(rolEf(ctx)) {
		ctx.JSON(http.StatusForbidden, gin.H{"error": "La comprobación en tiempo real es sólo para administrador y agente."})
		return
	}
	reportes, err := c.repo.UltimoCambioReportes(ctx.Request.Context())
	if err != nil {
		responder(ctx, err)
		return
	}
	msg, err := c.repo.UltimoMensajeID(ctx.Request.Context())
	if err != nil {
		responder(ctx, err)
		return
	}
	todos := ctx.MustGet("soporte_todos").([]entities.SoporteParticipante)
	ctx.JSON(http.StatusOK, gin.H{
		"reportes_actualizado": reportes,
		"ultimo_mensaje_id":    msg,
		"rol_efectivo":         rolEf(ctx),
		"participantes":        c.estado(todos, time.Now().UTC()),
	})
}

// ── reportes ───────────────────────────────────────────────────

func (c *SoporteController) CrearReporte(ctx *gin.Context) {
	var req soporte.NuevoReporte
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Cuerpo de la petición inválido."})
		return
	}
	in, err := soporte.ValidarNuevoReporte(req)
	if err != nil {
		responder(ctx, err)
		return
	}
	rc := ctx.Request.Context()

	// TC03 — política explícita: un ID repetido NUNCA sobrescribe; se responde
	// 409 y el reporte original queda intacto.
	if in.ID != "" {
		existe, err := c.repo.ExisteReporte(rc, in.ID)
		if err != nil {
			responder(ctx, err)
			return
		}
		if existe {
			ctx.JSON(http.StatusConflict, gin.H{"error": "Ya existe un reporte con el ID " + in.ID + ". No se sobrescribió."})
			return
		}
	}

	rep := &entities.SoporteReporte{
		ID: in.ID, Equipo: in.Equipo, Descripcion: in.Descripcion, Tipo: in.Tipo,
		Prioridad: in.Prioridad, Estado: soporte.EstadoAbierto, SolicitanteID: yo(ctx).ID,
		DirectoAdmin: in.DirectoAdmin, Activo: true,
	}
	for intento := 0; intento < 3; intento++ {
		if in.ID == "" {
			mayor, err := c.repo.MayorNumero(rc)
			if err != nil {
				responder(ctx, err)
				return
			}
			rep.ID = soporte.SiguienteID(mayor)
		}
		err = c.repo.CrearReporte(rc, rep)
		if err == nil {
			ctx.JSON(http.StatusCreated, rep)
			return
		}
		if !repositories.EsDuplicado(err) {
			responder(ctx, err)
			return
		}
		if in.ID != "" { // alguien lo creó entre la consulta y el INSERT
			ctx.JSON(http.StatusConflict, gin.H{"error": "Ya existe un reporte con el ID " + in.ID + ". No se sobrescribió."})
			return
		}
	}
	responder(ctx, err)
}

func (c *SoporteController) ListarReportes(ctx *gin.Context) {
	reps, err := c.repo.Reportes(ctx.Request.Context(), ctx.Query("estado"), ctx.Query("q"))
	if err != nil {
		responder(ctx, err)
		return
	}
	me := yo(ctx).ID
	out := []entities.SoporteReporte{}
	for _, r := range reps {
		if soporte.PuedeVer(rolEf(ctx), me, r) {
			out = append(out, r)
		}
	}
	ctx.JSON(http.StatusOK, out)
}

// reporteVisible responde 404 tanto si no existe como si el rol no puede
// verlo: así no se revela qué IDs existen (TC06).
func (c *SoporteController) reporteVisible(ctx *gin.Context) (*entities.SoporteReporte, string, bool) {
	id, ok := soporte.NormalizarID(ctx.Param("id"))
	noExiste := func() {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No existe un reporte activo con el ID " + id + "."})
	}
	if !ok {
		noExiste()
		return nil, id, false
	}
	rep, err := c.repo.Reporte(ctx.Request.Context(), id)
	if errors.Is(err, gorm.ErrRecordNotFound) || (err == nil && !soporte.PuedeVer(rolEf(ctx), yo(ctx).ID, *rep)) {
		noExiste()
		return nil, id, false
	}
	if err != nil {
		responder(ctx, err)
		return nil, id, false
	}
	return rep, id, true
}

func (c *SoporteController) VerReporte(ctx *gin.Context) {
	if rep, _, ok := c.reporteVisible(ctx); ok {
		ctx.JSON(http.StatusOK, rep)
	}
}

func (c *SoporteController) ActualizarReporte(ctx *gin.Context) {
	rep, id, ok := c.reporteVisible(ctx)
	if !ok {
		return
	}
	var cambio soporte.CambioReporte
	if err := ctx.ShouldBindJSON(&cambio); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Cuerpo de la petición inválido."})
		return
	}
	cols, err := soporte.ValidarCambio(rolEf(ctx), yo(ctx).ID, *rep, cambio)
	if err != nil {
		responder(ctx, err)
		return
	}
	if err := c.repo.ActualizarReporte(ctx.Request.Context(), id, cols); err != nil {
		responder(ctx, err)
		return
	}
	actualizado, err := c.repo.Reporte(ctx.Request.Context(), id)
	if err != nil {
		responder(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, actualizado)
}

// EliminarReporte es una baja lógica (activo = false): el registro deja de
// aparecer pero queda en la base para auditoría (TC05).
func (c *SoporteController) EliminarReporte(ctx *gin.Context) {
	if !soporte.PuedeEliminar(rolEf(ctx)) {
		responder(ctx, soporte.ErrSinPermiso)
		return
	}
	_, id, ok := c.reporteVisible(ctx)
	if !ok {
		return
	}
	if err := c.repo.ActualizarReporte(ctx.Request.Context(), id, map[string]any{"activo": false}); err != nil {
		responder(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// Asignar: el administrador (o el agente que lo cubre) asigna el reporte a un
// agente, o se lo queda él.
func (c *SoporteController) Asignar(ctx *gin.Context) {
	if !soporte.PuedeAdministrar(rolEf(ctx)) {
		responder(ctx, soporte.ErrSinPermiso)
		return
	}
	_, id, ok := c.reporteVisible(ctx)
	if !ok {
		return
	}
	var req struct {
		AgenteID string `json:"agente_id"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Cuerpo de la petición inválido."})
		return
	}
	destino, err := c.repo.Participante(ctx.Request.Context(), req.AgenteID)
	if err != nil || !destino.Activo || destino.Rol == soporte.RolUsuario {
		responder(ctx, &soporte.ErrValidacion{Campos: map[string]string{"agente_id": "debe ser un agente o administrador activo"}})
		return
	}
	if err := c.repo.ActualizarReporte(ctx.Request.Context(), id, map[string]any{
		"agente_id": destino.ID, "estado": soporte.EstadoAsignado,
	}); err != nil {
		responder(ctx, err)
		return
	}
	rep, err := c.repo.Reporte(ctx.Request.Context(), id)
	if err != nil {
		responder(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, rep)
}

// AsumirDesconectados: los reportes de un agente que perdió la conexión pasan
// a quien administra en ese momento.
func (c *SoporteController) AsumirDesconectados(ctx *gin.Context) {
	if !soporte.PuedeAdministrar(rolEf(ctx)) {
		responder(ctx, soporte.ErrSinPermiso)
		return
	}
	limite := time.Now().UTC().Add(-soporte.VentanaConexion)
	n, err := c.repo.AsumirDeDesconectados(ctx.Request.Context(), yo(ctx).ID, limite)
	if err != nil {
		responder(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"reasignados": n})
}

// ── mensajes y archivos ────────────────────────────────────────

func (c *SoporteController) ListarMensajes(ctx *gin.Context) {
	desde, _ := strconv.ParseInt(ctx.Query("desde"), 10, 64)
	me := yo(ctx).ID
	verTodo := rolEf(ctx) == soporte.RolAdmin
	ms, err := c.repo.Mensajes(ctx.Request.Context(), me, verTodo, desde)
	if err != nil {
		responder(ctx, err)
		return
	}
	conteo, err := c.repo.ConteoPorEmisor(ctx.Request.Context(), me, verTodo)
	if err != nil {
		responder(ctx, err)
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"mensajes": ms, "conteo": conteo})
}

type enviarMensajeRequest struct {
	Texto         string  `json:"texto"`
	ParaID        *string `json:"para_id"`
	ReporteID     *string `json:"reporte_id"`
	ArchivoPath   *string `json:"archivo_path"`
	ArchivoNombre *string `json:"archivo_nombre"`
}

func (c *SoporteController) EnviarMensaje(ctx *gin.Context) {
	var req enviarMensajeRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Cuerpo de la petición inválido."})
		return
	}
	me := yo(ctx)
	tieneArchivo := req.ArchivoPath != nil && *req.ArchivoPath != ""
	texto, err := soporte.ValidarMensaje(req.Texto, tieneArchivo)
	if err != nil {
		responder(ctx, err)
		return
	}
	if tieneArchivo && !soporte.EsRutaPropia(me.ID, *req.ArchivoPath) {
		responder(ctx, &soporte.ErrValidacion{Campos: map[string]string{"archivo": "sólo puedes adjuntar archivos que tú subiste"}})
		return
	}
	if req.ParaID != nil && *req.ParaID != "" {
		dest, err := c.repo.Participante(ctx.Request.Context(), *req.ParaID)
		if err != nil || !dest.Activo {
			responder(ctx, &soporte.ErrValidacion{Campos: map[string]string{"para_id": "el destinatario no existe"}})
			return
		}
	} else {
		req.ParaID = nil
	}
	if req.ReporteID != nil && *req.ReporteID != "" {
		id, ok := soporte.NormalizarID(*req.ReporteID)
		rep, err := c.repo.Reporte(ctx.Request.Context(), id)
		if !ok || err != nil || !soporte.PuedeVer(rolEf(ctx), me.ID, *rep) {
			responder(ctx, &soporte.ErrValidacion{Campos: map[string]string{"reporte_id": "el reporte no existe"}})
			return
		}
		req.ReporteID = &id
	} else {
		req.ReporteID = nil
	}
	if !tieneArchivo {
		req.ArchivoPath, req.ArchivoNombre = nil, nil
	}
	m := &entities.SoporteMensaje{
		ReporteID: req.ReporteID, EmisorID: me.ID, EmisorNombre: me.Nombre, EmisorRol: rolEf(ctx),
		ParaID: req.ParaID, Texto: texto, ArchivoPath: req.ArchivoPath, ArchivoNombre: req.ArchivoNombre,
		CreatedAt: time.Now().UTC(),
	}
	if err := c.repo.CrearMensaje(ctx.Request.Context(), m); err != nil {
		responder(ctx, err)
		return
	}
	ctx.JSON(http.StatusCreated, m)
}

// LimpiarMensajes borra el historial compartido; sólo quien administra.
// Los demás roles limpian su vista local en la app.
func (c *SoporteController) LimpiarMensajes(ctx *gin.Context) {
	if !soporte.PuedeAdministrar(rolEf(ctx)) {
		responder(ctx, soporte.ErrSinPermiso)
		return
	}
	if err := c.repo.BorrarMensajes(ctx.Request.Context()); err != nil {
		responder(ctx, err)
		return
	}
	ctx.Status(http.StatusNoContent)
}

// FirmarSubida entrega una URL preasignada de Supabase Storage para subir un
// archivo al bucket privado.
func (c *SoporteController) FirmarSubida(ctx *gin.Context) {
	var req struct {
		Nombre string `json:"nombre"`
		Tamano int64  `json:"tamano"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Cuerpo de la petición inválido."})
		return
	}
	seguro, mime, err := soporte.ValidarArchivo(req.Nombre, req.Tamano)
	if err != nil {
		responder(ctx, err)
		return
	}
	ruta := soporte.RutaArchivo(yo(ctx).ID, nuevoID(6), seguro)
	urlSubida, err := database.FirmarSubida(bucketSoporte, ruta)
	if err != nil {
		log.Printf("[soporte] firmar subida: %v", err)
		ctx.JSON(http.StatusBadGateway, gin.H{"error": "No fue posible preparar la subida del archivo."})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"upload_url": urlSubida, "path": ruta, "nombre": seguro, "content_type": mime})
}

// FirmarDescarga entrega una URL de lectura de 5 minutos, sólo si el archivo
// aparece en un mensaje que el solicitante puede ver.
func (c *SoporteController) FirmarDescarga(ctx *gin.Context) {
	var req struct {
		Path string `json:"path"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil || req.Path == "" {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Cuerpo de la petición inválido."})
		return
	}
	visible, err := c.repo.ArchivoVisible(ctx.Request.Context(), req.Path, yo(ctx).ID, rolEf(ctx) == soporte.RolAdmin)
	if err != nil {
		responder(ctx, err)
		return
	}
	if !visible {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "No se encontró el archivo."})
		return
	}
	u, err := database.FirmarDescarga(bucketSoporte, req.Path, 300)
	if err != nil {
		log.Printf("[soporte] firmar descarga: %v", err)
		ctx.JSON(http.StatusBadGateway, gin.H{"error": "No fue posible preparar la descarga del archivo."})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"url": u, "expira_en_segundos": 300})
}
