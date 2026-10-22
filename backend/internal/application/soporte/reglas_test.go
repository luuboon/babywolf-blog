package soporte

import (
	"errors"
	"strings"
	"testing"
	"time"

	"babywolf-blog-backend/internal/domain/entities"
)

var ahora = time.Date(2026, 10, 22, 12, 0, 0, 0, time.UTC)

func part(id, rol string, hace time.Duration) entities.SoporteParticipante {
	t := ahora.Add(-hace)
	return entities.SoporteParticipante{ID: id, Rol: rol, Activo: true, UltimoLatido: &t}
}

func ptr(s string) *string { return &s }

// TC02 — campos vacíos.
func TestReporteSinEquipoNiDescripcion(t *testing.T) {
	_, err := ValidarNuevoReporte(NuevoReporte{Equipo: "  ", Descripcion: ""})
	var v *ErrValidacion
	if !errors.As(err, &v) || v.Campos["equipo"] == "" || v.Campos["descripcion"] == "" {
		t.Fatalf("se esperaba validación de equipo y descripción, llegó %v", err)
	}
}

// TC01 — alta válida con valores por defecto.
func TestReporteValido(t *testing.T) {
	r, err := ValidarNuevoReporte(NuevoReporte{ID: " rep001 ", Equipo: "Laptop 3", Descripcion: "No enciende"})
	if err != nil {
		t.Fatal(err)
	}
	if r.ID != "REP001" || r.Prioridad != "Media" || r.Tipo != "falla" {
		t.Fatalf("normalización: %+v", r)
	}
}

func TestIDInvalido(t *testing.T) {
	for _, id := range []string{"REP1", "R001", "REP001X", "REP"} {
		if _, err := ValidarNuevoReporte(NuevoReporte{ID: id, Equipo: "a", Descripcion: "b"}); err == nil {
			t.Errorf("%q debió rechazarse", id)
		}
	}
}

func TestAccesoPerdidoVaDirectoAlAdmin(t *testing.T) {
	r, _ := ValidarNuevoReporte(NuevoReporte{Equipo: "Cuenta", Descripcion: "Olvidé mi contraseña", Tipo: "acceso"})
	if !r.DirectoAdmin {
		t.Fatal("una solicitud de acceso perdido debe ir directo al administrador")
	}
}

func TestSiguienteID(t *testing.T) {
	if SiguienteID(0) != "REP001" || SiguienteID(41) != "REP042" {
		t.Fatal("numeración REPnnn")
	}
}

// Intercambio de roles por desconexión.
func TestRolEfectivo(t *testing.T) {
	admin := part("ad", RolAdmin, 5*time.Second)
	agente := part("ag", RolAgente, 5*time.Second)
	if RolEfectivo(agente, []entities.SoporteParticipante{admin, agente}, ahora) != RolAgente {
		t.Fatal("con el admin conectado el agente sigue siendo agente")
	}
	adminFuera := part("ad", RolAdmin, 2*time.Minute)
	if RolEfectivo(agente, []entities.SoporteParticipante{adminFuera, agente}, ahora) != RolAdminTemporal {
		t.Fatal("sin admin conectado el agente cubre como admin temporal")
	}
	usuario := part("u1", RolUsuario, time.Second)
	if RolEfectivo(usuario, []entities.SoporteParticipante{usuario}, ahora) != RolUsuario {
		t.Fatal("un usuario nunca cambia de rol")
	}
}

func TestConectado(t *testing.T) {
	if !Conectado(part("x", RolAgente, 29*time.Second), ahora) {
		t.Fatal("29 s sigue conectado")
	}
	if Conectado(part("x", RolAgente, 31*time.Second), ahora) {
		t.Fatal("31 s ya está desconectado")
	}
	p := part("x", RolAgente, 0)
	p.UltimoLatido = nil
	if Conectado(p, ahora) {
		t.Fatal("tras Salir no hay latido")
	}
}

func TestPermisosDeVista(t *testing.T) {
	r := entities.SoporteReporte{SolicitanteID: "u1", AgenteID: ptr("ag")}
	if !PuedeVer(RolUsuario, "u1", r) || PuedeVer(RolUsuario, "u2", r) {
		t.Fatal("el usuario sólo ve sus reportes")
	}
	if !PuedeVer(RolAgente, "ag", r) || PuedeVer(RolAgente, "ag2", r) {
		t.Fatal("el agente ve lo asignado a él")
	}
	directo := entities.SoporteReporte{SolicitanteID: "u1", DirectoAdmin: true}
	if PuedeVer(RolAgente, "ag", directo) || !PuedeVer(RolAdmin, "ad", directo) {
		t.Fatal("las solicitudes directas sólo las ve el admin")
	}
}

// TC04 — cambio de prioridad y quién puede hacerlo.
func TestCambioDePrioridad(t *testing.T) {
	r := entities.SoporteReporte{SolicitanteID: "u1", AgenteID: ptr("ag"), Estado: EstadoAsignado}
	cols, err := ValidarCambio(RolAgente, "ag", r, CambioReporte{Prioridad: ptr("Alta")})
	if err != nil || cols["prioridad"] != "Alta" {
		t.Fatalf("el agente asignado puede subir la prioridad: %v", err)
	}
	if _, err := ValidarCambio(RolUsuario, "u1", r, CambioReporte{Prioridad: ptr("Alta")}); !errors.Is(err, ErrSinPermiso) {
		t.Fatal("el usuario no cambia prioridades")
	}
	if _, err := ValidarCambio(RolAdmin, "ad", r, CambioReporte{Prioridad: ptr("Urgente")}); err == nil {
		t.Fatal("prioridad fuera del catálogo")
	}
	if _, err := ValidarCambio(RolAdmin, "ad", r, CambioReporte{}); err == nil {
		t.Fatal("un cambio vacío no es válido")
	}
}

func TestEliminarSoloAdminReal(t *testing.T) {
	if !PuedeEliminar(RolAdmin) || PuedeEliminar(RolAdminTemporal) || PuedeEliminar(RolAgente) || PuedeEliminar(RolUsuario) {
		t.Fatal("sólo el administrador real elimina")
	}
}

func TestMonitorSoloAdminYAgente(t *testing.T) {
	if !PuedeMonitorear(RolAdmin) || !PuedeMonitorear(RolAgente) || !PuedeMonitorear(RolAdminTemporal) || PuedeMonitorear(RolUsuario) {
		t.Fatal("tiempo real sólo para admin y agente")
	}
}

func TestMensajeVacio(t *testing.T) {
	if _, err := ValidarMensaje("   ", false); err == nil {
		t.Fatal("mensaje vacío aceptado")
	}
	if _, err := ValidarMensaje("", true); err != nil {
		t.Fatal("un adjunto sin texto sí es válido")
	}
	if _, err := ValidarMensaje(strings.Repeat("x", MaxMensaje+1), false); err == nil {
		t.Fatal("mensaje demasiado largo aceptado")
	}
}

func TestArchivo(t *testing.T) {
	seguro, mime, err := ValidarArchivo("../../Captura de pantalla.PNG", 1000)
	if err != nil || seguro != "Captura_de_pantalla.png" || mime != "image/png" {
		t.Fatalf("nombre seguro: %q %q %v", seguro, mime, err)
	}
	if _, _, err := ValidarArchivo("virus.exe", 10); err == nil {
		t.Fatal("extensión no permitida aceptada")
	}
	if _, _, err := ValidarArchivo("foto.jpg", MaxArchivoBytes+1); err == nil {
		t.Fatal("archivo de más de 5 MB aceptado")
	}
	if !EsRutaPropia("p1", "soporte/p1/abc-foto.jpg") || EsRutaPropia("p1", "soporte/p2/abc-foto.jpg") || EsRutaPropia("p1", "soporte/p1/../p2/x.jpg") {
		t.Fatal("sólo se adjuntan archivos propios")
	}
}
