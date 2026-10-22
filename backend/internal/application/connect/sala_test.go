package connect

import (
	"errors"
	"strings"
	"testing"
)

func TestEnviarYListar(t *testing.T) {
	s := NewStore()
	if _, err := s.Enviar("lab-1", "a", "Pixel A", "  hola  "); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Enviar("LAB-1", "B", "", "qué tal"); err != nil {
		t.Fatal(err)
	}
	msgs, conteo, err := s.Listar("Lab-1", 0)
	if err != nil {
		t.Fatal(err)
	}
	if len(msgs) != 2 || msgs[0].Texto != "hola" || msgs[1].Dispositivo != "Dispositivo B" {
		t.Fatalf("mensajes inesperados: %+v", msgs)
	}
	if conteo["A"] != 1 || conteo["B"] != 1 {
		t.Fatalf("contador por dispositivo: %v", conteo)
	}
	nuevos, _, _ := s.Listar("LAB-1", msgs[0].ID)
	if len(nuevos) != 1 {
		t.Fatalf("desde debe filtrar lo ya visto: %d", len(nuevos))
	}
}

func TestValidaciones(t *testing.T) {
	s := NewStore()
	casos := []struct {
		sala, emisor, texto string
		esperado            error
	}{
		{"LAB-1", "A", "   ", ErrVacio},
		{"LAB-1", "C", "hola", ErrEmisor},
		{"x", "A", "hola", ErrSala},
		{"LAB 1", "A", "hola", ErrSala},
		{"LAB-1", "A", strings.Repeat("a", 281), ErrLargo},
	}
	for _, c := range casos {
		if _, err := s.Enviar(c.sala, c.emisor, "", c.texto); !errors.Is(err, c.esperado) {
			t.Errorf("%q/%q: se esperaba %v, llegó %v", c.sala, c.emisor, c.esperado, err)
		}
	}
}

func TestLimpiar(t *testing.T) {
	s := NewStore()
	s.Enviar("LAB-1", "A", "", "uno")
	if err := s.Limpiar("lab-1"); err != nil {
		t.Fatal(err)
	}
	msgs, conteo, _ := s.Listar("LAB-1", 0)
	if len(msgs) != 0 || conteo["A"] != 0 {
		t.Fatal("limpiar debe borrar mensajes y contadores")
	}
}

func TestTopePorSala(t *testing.T) {
	s := NewStore()
	for i := 0; i < maxPorSala+10; i++ {
		s.Enviar("LAB-1", "A", "", "m")
	}
	msgs, conteo, _ := s.Listar("LAB-1", 0)
	if len(msgs) != maxPorSala {
		t.Fatalf("la sala debe guardar sólo %d mensajes, guarda %d", maxPorSala, len(msgs))
	}
	if conteo["A"] != maxPorSala+10 {
		t.Fatalf("el contador cuenta todos los enviados: %d", conteo["A"])
	}
}
