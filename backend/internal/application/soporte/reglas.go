// Package soporte contiene las reglas de negocio de la mesa de ayuda
// (Práctica 7). Son funciones puras, sin base de datos, para poder probarlas
// con go test y citarlas en la matriz de caja negra (Práctica 8).
package soporte

import (
	"errors"
	"fmt"
	"path"
	"regexp"
	"strings"
	"time"
	"unicode/utf8"

	"babywolf-blog-backend/internal/domain/entities"
)

const (
	RolAdmin         = "admin"
	RolAgente        = "agente"
	RolUsuario       = "usuario"
	RolAdminTemporal = "admin_temporal" // agente que cubre al admin desconectado

	EstadoAbierto   = "abierto"
	EstadoAsignado  = "asignado"
	EstadoProceso   = "en_proceso"
	EstadoResuelto  = "resuelto"
	MaxEquipo       = 80
	MaxDescripcion  = 500
	MaxMensaje      = 500
	MaxNombre       = 40
	MaxArchivoBytes = 5 * 1024 * 1024

	// Un dispositivo cuenta como conectado si mandó un latido en los
	// últimos 30 s. La app manda uno cada 10 s.
	VentanaConexion = 30 * time.Second
)

var (
	tipos       = map[string]bool{"falla": true, "atencion": true, "acceso": true}
	prioridades = map[string]bool{"Baja": true, "Media": true, "Alta": true}
	estados     = map[string]bool{EstadoAbierto: true, EstadoAsignado: true, EstadoProceso: true, EstadoResuelto: true}
	reID        = regexp.MustCompile(`^REP[0-9]{3,6}$`)
	reArchivo   = regexp.MustCompile(`[^a-zA-Z0-9._-]+`)
	extensiones = map[string]string{
		".png": "image/png", ".jpg": "image/jpeg", ".jpeg": "image/jpeg",
		".pdf": "application/pdf", ".txt": "text/plain", ".log": "text/plain",
	}

	ErrSinPermiso = errors.New("no tienes permiso para esta acción")
)

// ErrValidacion junta los campos inválidos para que la app los muestre
// junto a cada campo (TC02).
type ErrValidacion struct {
	Campos map[string]string
}

func (e *ErrValidacion) Error() string {
	partes := []string{}
	for k, v := range e.Campos {
		partes = append(partes, k+": "+v)
	}
	return "datos inválidos (" + strings.Join(partes, "; ") + ")"
}

func Conectado(p entities.SoporteParticipante, ahora time.Time) bool {
	return p.Activo && p.UltimoLatido != nil && ahora.Sub(*p.UltimoLatido) <= VentanaConexion
}

// RolEfectivo implementa el "intercambio de roles por desconexión": si no hay
// ningún administrador conectado, el agente asume funciones de administrador
// (asignar, invitar, reasignar) hasta que el admin vuelva.
func RolEfectivo(p entities.SoporteParticipante, todos []entities.SoporteParticipante, ahora time.Time) string {
	if p.Rol != RolAgente {
		return p.Rol
	}
	for _, o := range todos {
		if o.Rol == RolAdmin && Conectado(o, ahora) {
			return RolAgente
		}
	}
	return RolAdminTemporal
}

func PuedeAdministrar(rolEfectivo string) bool {
	return rolEfectivo == RolAdmin || rolEfectivo == RolAdminTemporal
}

// PuedeMonitorear: la comprobación en tiempo real es sólo para admin y agente.
func PuedeMonitorear(rolEfectivo string) bool {
	return rolEfectivo != RolUsuario
}

func ValidarRolInvitacion(rol string) error {
	if rol != RolAgente && rol != RolUsuario {
		return &ErrValidacion{Campos: map[string]string{"rol": "sólo se invita como agente o usuario"}}
	}
	return nil
}

func ValidarNombre(nombre string) (string, error) {
	n := strings.TrimSpace(nombre)
	if n == "" {
		return "", &ErrValidacion{Campos: map[string]string{"nombre": "es obligatorio"}}
	}
	if utf8.RuneCountInString(n) > MaxNombre {
		return "", &ErrValidacion{Campos: map[string]string{"nombre": fmt.Sprintf("máximo %d caracteres", MaxNombre)}}
	}
	return n, nil
}

type NuevoReporte struct {
	ID           string `json:"id"`
	Equipo       string `json:"equipo"`
	Descripcion  string `json:"descripcion"`
	Tipo         string `json:"tipo"`
	Prioridad    string `json:"prioridad"`
	DirectoAdmin bool   `json:"directo_admin"`
}

// NormalizarID acepta "rep001" o " REP001 " y devuelve "REP001".
func NormalizarID(id string) (string, bool) {
	n := strings.ToUpper(strings.TrimSpace(id))
	return n, reID.MatchString(n)
}

// ValidarNuevoReporte aplica TC02: equipo y descripción obligatorios.
func ValidarNuevoReporte(in NuevoReporte) (NuevoReporte, error) {
	campos := map[string]string{}
	in.Equipo = strings.TrimSpace(in.Equipo)
	in.Descripcion = strings.TrimSpace(in.Descripcion)
	if in.Equipo == "" {
		campos["equipo"] = "es obligatorio"
	} else if utf8.RuneCountInString(in.Equipo) > MaxEquipo {
		campos["equipo"] = fmt.Sprintf("máximo %d caracteres", MaxEquipo)
	}
	if in.Descripcion == "" {
		campos["descripcion"] = "es obligatoria"
	} else if utf8.RuneCountInString(in.Descripcion) > MaxDescripcion {
		campos["descripcion"] = fmt.Sprintf("máximo %d caracteres", MaxDescripcion)
	}
	if in.Tipo == "" {
		in.Tipo = "falla"
	}
	if !tipos[in.Tipo] {
		campos["tipo"] = "debe ser falla, atencion o acceso"
	}
	if in.Prioridad == "" {
		in.Prioridad = "Media"
	}
	if !prioridades[in.Prioridad] {
		campos["prioridad"] = "debe ser Baja, Media o Alta"
	}
	if strings.TrimSpace(in.ID) != "" {
		id, ok := NormalizarID(in.ID)
		if !ok {
			campos["id"] = "formato REP seguido de 3 a 6 dígitos, p. ej. REP001"
		}
		in.ID = id
	}
	// Las solicitudes de acceso perdido van siempre directo al administrador.
	if in.Tipo == "acceso" {
		in.DirectoAdmin = true
	}
	if len(campos) > 0 {
		return in, &ErrValidacion{Campos: campos}
	}
	return in, nil
}

// SiguienteID numera a partir del mayor existente (REP001, REP002…).
func SiguienteID(mayor int) string {
	return fmt.Sprintf("REP%03d", mayor+1)
}

// PuedeVer: admin ve todo; el agente ve lo asignado a él y lo abierto sin
// asignar que no sea solicitud directa al admin; el usuario sólo lo suyo.
func PuedeVer(rolEf, yo string, r entities.SoporteReporte) bool {
	switch {
	case PuedeAdministrar(rolEf):
		return true
	case rolEf == RolAgente:
		return (r.AgenteID != nil && *r.AgenteID == yo) || (r.AgenteID == nil && !r.DirectoAdmin)
	default:
		return r.SolicitanteID == yo
	}
}

type CambioReporte struct {
	Equipo      *string `json:"equipo"`
	Descripcion *string `json:"descripcion"`
	Prioridad   *string `json:"prioridad"`
	Estado      *string `json:"estado"`
}

// ValidarCambio decide qué puede modificar cada rol (TC04) y devuelve las
// columnas a actualizar.
func ValidarCambio(rolEf, yo string, r entities.SoporteReporte, c CambioReporte) (map[string]any, error) {
	esAdmin := PuedeAdministrar(rolEf)
	esSuAgente := rolEf == RolAgente && r.AgenteID != nil && *r.AgenteID == yo
	esDueno := r.SolicitanteID == yo && r.Estado == EstadoAbierto

	cols := map[string]any{}
	campos := map[string]string{}

	if c.Prioridad != nil || c.Estado != nil {
		if !esAdmin && !esSuAgente {
			return nil, ErrSinPermiso
		}
		if c.Prioridad != nil {
			if !prioridades[*c.Prioridad] {
				campos["prioridad"] = "debe ser Baja, Media o Alta"
			}
			cols["prioridad"] = *c.Prioridad
		}
		if c.Estado != nil {
			if !estados[*c.Estado] {
				campos["estado"] = "estado desconocido"
			}
			cols["estado"] = *c.Estado
		}
	}
	if c.Equipo != nil || c.Descripcion != nil {
		if !esAdmin && !esDueno {
			return nil, ErrSinPermiso
		}
		if c.Equipo != nil {
			v := strings.TrimSpace(*c.Equipo)
			if v == "" || utf8.RuneCountInString(v) > MaxEquipo {
				campos["equipo"] = "obligatorio, máximo 80 caracteres"
			}
			cols["equipo"] = v
		}
		if c.Descripcion != nil {
			v := strings.TrimSpace(*c.Descripcion)
			if v == "" || utf8.RuneCountInString(v) > MaxDescripcion {
				campos["descripcion"] = "obligatoria, máximo 500 caracteres"
			}
			cols["descripcion"] = v
		}
	}
	if len(campos) > 0 {
		return nil, &ErrValidacion{Campos: campos}
	}
	if len(cols) == 0 {
		return nil, &ErrValidacion{Campos: map[string]string{"cambio": "no hay nada que actualizar"}}
	}
	return cols, nil
}

// PuedeEliminar: borrar (baja lógica) es exclusivo del administrador real;
// el agente que cubre al admin no puede eliminar.
func PuedeEliminar(rolEf string) bool {
	return rolEf == RolAdmin
}

// ValidarMensaje: no se aceptan mensajes vacíos salvo que lleven archivo.
func ValidarMensaje(texto string, tieneArchivo bool) (string, error) {
	t := strings.TrimSpace(texto)
	if t == "" && !tieneArchivo {
		return "", &ErrValidacion{Campos: map[string]string{"texto": "el mensaje no puede estar vacío"}}
	}
	if utf8.RuneCountInString(t) > MaxMensaje {
		return "", &ErrValidacion{Campos: map[string]string{"texto": fmt.Sprintf("máximo %d caracteres", MaxMensaje)}}
	}
	return t, nil
}

// ValidarArchivo revisa extensión y tamaño, y devuelve un nombre seguro y su
// tipo MIME. El nombre seguro evita rutas como "../../otro".
func ValidarArchivo(nombre string, tamano int64) (seguro, mime string, err error) {
	base := path.Base(strings.ReplaceAll(nombre, "\\", "/"))
	ext := strings.ToLower(path.Ext(base))
	mime, ok := extensiones[ext]
	if !ok {
		return "", "", &ErrValidacion{Campos: map[string]string{"archivo": "sólo png, jpg, pdf, txt o log"}}
	}
	if tamano <= 0 || tamano > MaxArchivoBytes {
		return "", "", &ErrValidacion{Campos: map[string]string{"archivo": "tamaño entre 1 byte y 5 MB"}}
	}
	seguro = reArchivo.ReplaceAllString(strings.TrimSuffix(base, path.Ext(base)), "_")
	seguro = strings.Trim(seguro, "._-")
	if seguro == "" {
		seguro = "archivo"
	}
	if utf8.RuneCountInString(seguro) > 60 {
		seguro = string([]rune(seguro)[:60])
	}
	return seguro + ext, mime, nil
}

// RutaArchivo: cada participante sólo puede adjuntar archivos de su carpeta.
func RutaArchivo(participanteID, id, seguro string) string {
	return "soporte/" + participanteID + "/" + id + "-" + seguro
}

func EsRutaPropia(participanteID, ruta string) bool {
	return strings.HasPrefix(ruta, "soporte/"+participanteID+"/") && !strings.Contains(ruta, "..")
}
