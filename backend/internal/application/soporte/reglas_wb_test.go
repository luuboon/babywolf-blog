package soporte

// Práctica 11 — caja blanca de ValidarCambio (el método con más decisiones del
// backend). Cada caso nombra el camino del grafo que ejecuta; ver
// docs/practicas/practica11_whitebox/metodos_criticos.md en el repo de la app.

import (
	"errors"
	"strings"
	"testing"

	"babywolf-blog-backend/internal/domain/entities"
)

func reporte(solicitante string, estado string, agente *string) entities.SoporteReporte {
	return entities.SoporteReporte{SolicitanteID: solicitante, Estado: estado, AgenteID: agente}
}

func TestWB_ValidarCambio_Caminos(t *testing.T) {
	yo := "u1"
	otro := "u2"
	abierto := reporte(yo, EstadoAbierto, nil)
	asignadoAYo := reporte("x", EstadoAsignado, &yo)
	asignadoAOtro := reporte("x", EstadoAsignado, &otro)
	proceso := reporte(yo, EstadoProceso, nil)

	casos := []struct {
		nombre string
		rol    string
		rep    entities.SoporteReporte
		cambio CambioReporte
		cols   []string // columnas esperadas si no hay error
		err    string   // "permiso", "validacion" o ""
		campo  string   // campo inválido esperado
	}{
		// P1: ningún campo → "no hay nada que actualizar"
		{"P1 sin campos", RolAdmin, abierto, CambioReporte{}, nil, "validacion", "cambio"},
		// P2: admin cambia prioridad válida
		{"P2 admin prioridad válida", RolAdmin, abierto, CambioReporte{Prioridad: ptr("Alta")}, []string{"prioridad"}, "", ""},
		// P3: usuario común intenta cambiar prioridad → sin permiso
		{"P3 usuario cambia prioridad", RolUsuario, abierto, CambioReporte{Prioridad: ptr("Alta")}, nil, "permiso", ""},
		// P4: prioridad inválida (admin)
		{"P4 prioridad inválida", RolAdmin, abierto, CambioReporte{Prioridad: ptr("Urgente")}, nil, "validacion", "prioridad"},
		// P5: agente asignado cambia estado válido
		{"P5 su agente cambia estado", RolAgente, asignadoAYo, CambioReporte{Estado: ptr(EstadoProceso)}, []string{"estado"}, "", ""},
		// P6: agente NO asignado cambia estado → sin permiso
		{"P6 otro agente cambia estado", RolAgente, asignadoAOtro, CambioReporte{Estado: ptr(EstadoProceso)}, nil, "permiso", ""},
		// P7: estado inválido
		{"P7 estado inválido", RolAdmin, abierto, CambioReporte{Estado: ptr("cerrado")}, nil, "validacion", "estado"},
		// P8: dueño edita equipo y descripción mientras está abierto
		{"P8 dueño edita abierto", RolUsuario, abierto, CambioReporte{Equipo: ptr(" PC-12 "), Descripcion: ptr("No enciende")}, []string{"equipo", "descripcion"}, "", ""},
		// P9: dueño edita cuando ya no está abierto → sin permiso
		{"P9 dueño edita en proceso", RolUsuario, proceso, CambioReporte{Equipo: ptr("PC")}, nil, "permiso", ""},
		// P10: admin con equipo vacío
		{"P10 equipo vacío", RolAdmin, abierto, CambioReporte{Equipo: ptr("   ")}, nil, "validacion", "equipo"},
		// P11: equipo de 81 caracteres (límite + 1)
		{"P11 equipo 81", RolAdmin, abierto, CambioReporte{Equipo: ptr(strings.Repeat("e", MaxEquipo+1))}, nil, "validacion", "equipo"},
		// P12: equipo de 80 caracteres (límite)
		{"P12 equipo 80", RolAdmin, abierto, CambioReporte{Equipo: ptr(strings.Repeat("e", MaxEquipo))}, []string{"equipo"}, "", ""},
		// P13: descripción vacía
		{"P13 descripción vacía", RolAdmin, abierto, CambioReporte{Descripcion: ptr("")}, nil, "validacion", "descripcion"},
		// P14: descripción de 501 y de 500
		{"P14 descripción 501", RolAdmin, abierto, CambioReporte{Descripcion: ptr(strings.Repeat("d", MaxDescripcion+1))}, nil, "validacion", "descripcion"},
		{"P14b descripción 500", RolAdmin, abierto, CambioReporte{Descripcion: ptr(strings.Repeat("d", MaxDescripcion))}, []string{"descripcion"}, "", ""},
		// P15: agente que cubre al admin (admin_temporal) puede cambiar todo
		{"P15 admin temporal", RolAdminTemporal, asignadoAOtro, CambioReporte{Prioridad: ptr("Baja"), Equipo: ptr("PC")}, []string{"prioridad", "equipo"}, "", ""},
		// P16: varios campos inválidos a la vez se reportan juntos
		{"P16 errores múltiples", RolAdmin, abierto, CambioReporte{Prioridad: ptr("x"), Estado: ptr("y"), Equipo: ptr("")}, nil, "validacion", "prioridad"},
	}

	for _, c := range casos {
		t.Run(c.nombre, func(t *testing.T) {
			cols, err := ValidarCambio(c.rol, yo, c.rep, c.cambio)
			switch c.err {
			case "permiso":
				if !errors.Is(err, ErrSinPermiso) {
					t.Fatalf("se esperaba ErrSinPermiso, llegó %v", err)
				}
			case "validacion":
				var v *ErrValidacion
				if !errors.As(err, &v) {
					t.Fatalf("se esperaba ErrValidacion, llegó %v", err)
				}
				if _, ok := v.Campos[c.campo]; !ok {
					t.Fatalf("falta el campo %q en %v", c.campo, v.Campos)
				}
			default:
				if err != nil {
					t.Fatalf("no se esperaba error: %v", err)
				}
				if len(cols) != len(c.cols) {
					t.Fatalf("columnas %v, se esperaban %v", cols, c.cols)
				}
				for _, k := range c.cols {
					if _, ok := cols[k]; !ok {
						t.Fatalf("falta la columna %q en %v", k, cols)
					}
				}
			}
		})
	}
}

// Se recorta antes de guardar: el valor que llega a la base no lleva espacios.
func TestWB_ValidarCambio_RecortaEspacios(t *testing.T) {
	cols, err := ValidarCambio(RolAdmin, "a", reporte("a", EstadoAbierto, nil), CambioReporte{Equipo: ptr("  PC-1  ")})
	if err != nil || cols["equipo"] != "PC-1" {
		t.Fatalf("equipo sin recortar: %v %v", cols, err)
	}
}

func TestWB_ErrValidacion_Error(t *testing.T) {
	e := &ErrValidacion{Campos: map[string]string{"a": "mal"}}
	if !strings.Contains(e.Error(), "a: mal") {
		t.Fatalf("mensaje: %q", e.Error())
	}
}

func TestWB_ValidarRolInvitacion_Y_Nombre(t *testing.T) {
	if ValidarRolInvitacion(RolAgente) != nil || ValidarRolInvitacion(RolUsuario) != nil {
		t.Fatal("agente y usuario se pueden invitar")
	}
	if ValidarRolInvitacion(RolAdmin) == nil {
		t.Fatal("no se invita como admin")
	}
	if _, err := ValidarNombre("  "); err == nil {
		t.Fatal("nombre vacío aceptado")
	}
	if _, err := ValidarNombre(strings.Repeat("n", MaxNombre+1)); err == nil {
		t.Fatal("nombre largo aceptado")
	}
	if n, err := ValidarNombre("  Ana "); err != nil || n != "Ana" {
		t.Fatalf("nombre válido: %q %v", n, err)
	}
}

func TestWB_ValidarArchivo_Limites(t *testing.T) {
	if _, _, err := ValidarArchivo("../../etc/passwd.txt", 10); err != nil {
		t.Fatalf("la ruta se reduce al nombre base: %v", err)
	}
	s, _, _ := ValidarArchivo("../../ra\\iz/foto mala.PNG", 10)
	if s != "foto_mala.png" {
		t.Fatalf("nombre seguro: %q", s)
	}
	if _, _, err := ValidarArchivo("virus.exe", 10); err == nil {
		t.Fatal("extensión no permitida aceptada")
	}
	if _, _, err := ValidarArchivo("a.png", 0); err == nil {
		t.Fatal("0 bytes aceptado")
	}
	if _, _, err := ValidarArchivo("a.png", MaxArchivoBytes+1); err == nil {
		t.Fatal("más de 5 MB aceptado")
	}
	if _, _, err := ValidarArchivo("a.png", MaxArchivoBytes); err != nil {
		t.Fatalf("exactamente 5 MB debe pasar: %v", err)
	}
	if s, _, _ := ValidarArchivo("....png", 5); s != "archivo.png" {
		t.Fatalf("nombre vacío tras limpiar: %q", s)
	}
	largo := strings.Repeat("x", 100) + ".txt"
	if s, _, _ := ValidarArchivo(largo, 5); len([]rune(s)) != 64 {
		t.Fatalf("nombre recortado a 60 + extensión: %q", s)
	}
}

func TestWB_Rutas(t *testing.T) {
	r := RutaArchivo("p1", "m9", "foto.png")
	if r != "soporte/p1/m9-foto.png" {
		t.Fatal(r)
	}
	if !EsRutaPropia("p1", r) || EsRutaPropia("p2", r) || EsRutaPropia("p1", "soporte/p1/../p2/x") {
		t.Fatal("EsRutaPropia")
	}
}

// ValidarNuevoReporte: ramas de longitud máxima, tipo y prioridad inválidos.
func TestWB_ValidarNuevoReporte_Ramas(t *testing.T) {
	_, err := ValidarNuevoReporte(NuevoReporte{
		Equipo:      strings.Repeat("e", MaxEquipo+1),
		Descripcion: strings.Repeat("d", MaxDescripcion+1),
		Tipo:        "otro",
		Prioridad:   "Urgente",
		ID:          "REP1",
	})
	var v *ErrValidacion
	if !errors.As(err, &v) {
		t.Fatalf("se esperaba ErrValidacion, llegó %v", err)
	}
	for _, k := range []string{"equipo", "descripcion", "tipo", "prioridad", "id"} {
		if v.Campos[k] == "" {
			t.Errorf("falta el error de %q: %v", k, v.Campos)
		}
	}
	// Límites exactos: 80 y 500 pasan.
	r, err := ValidarNuevoReporte(NuevoReporte{Equipo: strings.Repeat("e", MaxEquipo), Descripcion: strings.Repeat("d", MaxDescripcion), Tipo: "acceso", ID: " rep007 "})
	if err != nil || !r.DirectoAdmin || r.ID != "REP007" || r.Prioridad != "Media" {
		t.Fatalf("límites exactos: %+v %v", r, err)
	}
}
