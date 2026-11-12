package connect

// Práctica 11 — pruebas de caja blanca de connect.Store.
//
// Cubren las ramas que las pruebas funcionales (sala_test.go) no ejecutan:
// valores límite de cada condición, purga por inactividad y por tope de
// salas. Cada prueba indica el camino del grafo de flujo que ejecuta.

import (
	"errors"
	"strings"
	"testing"
	"time"
)

// WB-C1 · valores límite de NormalizarSala (condición: 3 <= largo <= 24).
func TestWB_NormalizarSala_Limites(t *testing.T) {
	casos := []struct {
		entrada string
		ok      bool
	}{
		{"AB", false},                    // largo 2: justo debajo del mínimo
		{"ABC", true},                    // largo 3: mínimo
		{strings.Repeat("A", 24), true},  // largo 24: máximo
		{strings.Repeat("A", 25), false}, // largo 25: justo arriba
		{"  lab-1  ", true},              // recorta espacios y pasa a mayúsculas
		{"lab_1", false},                 // guion bajo no permitido
		{"", false},                      // vacío
	}
	for _, c := range casos {
		_, err := NormalizarSala(c.entrada)
		if (err == nil) != c.ok {
			t.Errorf("NormalizarSala(%q): ok=%v, error=%v", c.entrada, c.ok, err)
		}
	}
}

// WB-C2 · límites del texto (280) contado en runas, no en bytes.
func TestWB_Enviar_LimiteDeTextoEnRunas(t *testing.T) {
	s := NewStore()
	// 280 caracteres de 2 bytes cada uno = 560 bytes: debe pasar.
	if _, err := s.Enviar("LAB-1", "A", "", strings.Repeat("ñ", MaxTexto)); err != nil {
		t.Fatalf("280 runas deben pasar: %v", err)
	}
	if _, err := s.Enviar("LAB-1", "A", "", strings.Repeat("ñ", MaxTexto+1)); !errors.Is(err, ErrLargo) {
		t.Fatalf("281 runas deben fallar con ErrLargo, llegó %v", err)
	}
}

// WB-C3 · nombre del dispositivo: 30 pasa, 31 falla, vacío usa el valor por omisión.
func TestWB_Enviar_Dispositivo(t *testing.T) {
	s := NewStore()
	if _, err := s.Enviar("LAB-1", "A", strings.Repeat("d", MaxDispositivo), "x"); err != nil {
		t.Fatalf("30 caracteres deben pasar: %v", err)
	}
	if _, err := s.Enviar("LAB-1", "A", strings.Repeat("d", MaxDispositivo+1), "x"); !errors.Is(err, ErrDispositivo) {
		t.Fatalf("31 caracteres deben fallar con ErrDispositivo, llegó %v", err)
	}
	m, _ := s.Enviar("LAB-1", "b", "   ", "x")
	if m.Dispositivo != "Dispositivo B" || m.Emisor != "B" {
		t.Fatalf("valores por omisión: %+v", m)
	}
}

// WB-C4 · purgar: una sala sin actividad por más de 24 h se descarta.
func TestWB_Purgar_PorInactividad(t *testing.T) {
	s := NewStore()
	ahora := time.Date(2026, 11, 11, 12, 0, 0, 0, time.UTC)
	s.Ahora = func() time.Time { return ahora }
	s.Enviar("VIEJA", "A", "", "hola")
	// Justo en el límite (24 h) NO se purga: la condición es estrictamente mayor.
	s.Ahora = func() time.Time { return ahora.Add(vidaSala) }
	s.Enviar("NUEVA", "A", "", "hola")
	if _, ok := s.salas["VIEJA"]; !ok {
		t.Fatal("a las 24 h exactas la sala todavía debe existir")
	}
	// Un segundo después sí.
	s.Ahora = func() time.Time { return ahora.Add(vidaSala + time.Second) }
	s.Enviar("OTRA", "A", "", "hola")
	if _, ok := s.salas["VIEJA"]; ok {
		t.Fatal("pasadas 24 h la sala debe purgarse")
	}
}

// WB-C5 · purgar: más de maxSalas descarta la menos activa.
func TestWB_Purgar_TopeDeSalas(t *testing.T) {
	s := NewStore()
	base := time.Date(2026, 11, 11, 12, 0, 0, 0, time.UTC)
	n := 0
	s.Ahora = func() time.Time { n++; return base.Add(time.Duration(n) * time.Second) }
	for i := 0; i <= maxSalas; i++ { // maxSalas+1 salas distintas
		codigo := "SALA-" + strings.Repeat("A", 3) + string(rune('A'+i%26)) + string(rune('A'+(i/26)%26))
		if _, err := s.Enviar(codigo, "A", "", "x"); err != nil {
			t.Fatal(err)
		}
	}
	// La purga corre al INICIO de Enviar, así que tras la sala 201 todavía hay 201.
	if len(s.salas) != maxSalas+1 {
		t.Fatalf("antes de la siguiente purga hay %d salas", len(s.salas))
	}
	s.Enviar("SALA-ZZZZZ", "A", "", "x") // dispara la purga: se va la más antigua
	if len(s.salas) > maxSalas+1 {
		t.Fatalf("el tope no se respeta: %d salas", len(s.salas))
	}
	if _, ok := s.salas["SALA-AAAAA"]; ok {
		t.Fatal("la sala menos activa debía descartarse")
	}
}

// WB-C6 · Listar y Limpiar con código inválido o sala inexistente.
func TestWB_Listar_Limpiar_Excepciones(t *testing.T) {
	s := NewStore()
	if _, _, err := s.Listar("x", 0); !errors.Is(err, ErrSala) {
		t.Fatalf("Listar con código inválido: %v", err)
	}
	if err := s.Limpiar("x"); !errors.Is(err, ErrSala) {
		t.Fatalf("Limpiar con código inválido: %v", err)
	}
	msgs, conteo, err := s.Listar("NOEXISTE", 0)
	if err != nil || len(msgs) != 0 || conteo["A"] != 0 || conteo["B"] != 0 {
		t.Fatalf("sala inexistente debe dar lista vacía: %v %v %v", msgs, conteo, err)
	}
}

// WB-C7 · concurrencia: 50 goroutines escribiendo en la misma sala no pierden
// mensajes ni repiten ids (se corre con -race).
func TestWB_Enviar_Concurrente(t *testing.T) {
	s := NewStore()
	hecho := make(chan struct{})
	const n = 50
	for i := 0; i < n; i++ {
		go func() {
			defer func() { hecho <- struct{}{} }()
			s.Enviar("LAB-1", "A", "", "m")
		}()
	}
	for i := 0; i < n; i++ {
		<-hecho
	}
	msgs, conteo, _ := s.Listar("LAB-1", 0)
	if len(msgs) != n || conteo["A"] != n {
		t.Fatalf("se esperaban %d mensajes, hay %d (contador %d)", n, len(msgs), conteo["A"])
	}
	vistos := map[int64]bool{}
	for _, m := range msgs {
		if vistos[m.ID] {
			t.Fatalf("id repetido: %d", m.ID)
		}
		vistos[m.ID] = true
	}
}
