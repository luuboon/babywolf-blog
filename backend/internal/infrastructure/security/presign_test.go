package security

import (
	"errors"
	"net/url"
	"testing"
	"time"
)

func init() {
	ClaveFn = func() ([]byte, error) { return []byte("clave-de-prueba"), nil }
}

var t0 = time.Date(2026, 10, 22, 12, 0, 0, 0, time.UTC)

func TestFirmarYVerificar(t *testing.T) {
	q, _, err := Firmar("consulta", url.Values{"recurso": {"resumen"}}, 5*time.Minute, t0)
	if err != nil {
		t.Fatal(err)
	}
	if err := Verificar("consulta", q, t0.Add(time.Minute)); err != nil {
		t.Fatalf("una URL íntegra debe pasar: %v", err)
	}
}

// Prueba 3 de la práctica 3: modificar la URL invalida la firma.
func TestURLModificadaSeRechaza(t *testing.T) {
	q, _, _ := Firmar("consulta", url.Values{"recurso": {"resumen"}}, 5*time.Minute, t0)
	casos := map[string]func(url.Values){
		"parametro cambiado": func(v url.Values) { v.Set("recurso", "otro") },
		"parametro agregado": func(v url.Values) { v.Set("admin", "1") },
		"exp alargado":       func(v url.Values) { v.Set("exp", "99999999999") },
		"firma alterada":     func(v url.Values) { v.Set("sig", v.Get("sig")[1:]+"A") },
		"sin firma":          func(v url.Values) { v.Del("sig") },
	}
	for nombre, alterar := range casos {
		v, _ := url.ParseQuery(q.Encode())
		alterar(v)
		if err := Verificar("consulta", v, t0); !errors.Is(err, ErrFirmaInvalida) {
			t.Errorf("%s: se esperaba ErrFirmaInvalida, llegó %v", nombre, err)
		}
	}
}

func TestURLExpirada(t *testing.T) {
	q, _, _ := Firmar("consulta", url.Values{}, 5*time.Minute, t0)
	if err := Verificar("consulta", q, t0.Add(6*time.Minute)); !errors.Is(err, ErrExpirada) {
		t.Fatalf("se esperaba ErrExpirada, llegó %v", err)
	}
}

func TestPropositoSeparaFirmas(t *testing.T) {
	q, _, _ := Firmar("consulta", url.Values{"rol": {"agente"}}, time.Minute, t0)
	if err := Verificar("soporte-invitacion", q, t0); !errors.Is(err, ErrFirmaInvalida) {
		t.Fatal("una firma de consulta no debe servir como invitación")
	}
}

func TestToken(t *testing.T) {
	tok, err := EmitirToken("abc-123", time.Hour, t0)
	if err != nil {
		t.Fatal(err)
	}
	id, err := LeerToken(tok, t0.Add(30*time.Minute))
	if err != nil || id != "abc-123" {
		t.Fatalf("token válido rechazado: %v %q", err, id)
	}
	if _, err := LeerToken("otro-id"+tok[len("abc-123"):], t0); !errors.Is(err, ErrFirmaInvalida) {
		t.Fatal("cambiar el id del token debe invalidarlo")
	}
	if _, err := LeerToken(tok, t0.Add(2*time.Hour)); !errors.Is(err, ErrExpirada) {
		t.Fatal("el token debe expirar")
	}
	if _, err := LeerToken("", t0); err == nil {
		t.Fatal("token vacío aceptado")
	}
}
