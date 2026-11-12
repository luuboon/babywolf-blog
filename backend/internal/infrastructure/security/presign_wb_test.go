package security

// Práctica 11 — caja blanca de presign.go: ramas de claveDesdeEntorno y de
// manejo de errores que las pruebas funcionales no recorren.

import (
	"errors"
	"net/url"
	"strconv"
	"testing"
	"time"
)

func sinClave() { ClaveFn = func() ([]byte, error) { return nil, ErrSinClave } }
func conClave() { ClaveFn = func() ([]byte, error) { return []byte("clave-de-prueba"), nil } }

// Orden de preferencia de la clave: PRESIGN_SECRET > JWT_SECRET > SERVICE_ROLE.
func TestWB_ClaveDesdeEntorno_Ramas(t *testing.T) {
	t.Setenv("PRESIGN_SECRET", "")
	t.Setenv("JWT_SECRET", "")
	t.Setenv("SUPABASE_SERVICE_ROLE_KEY", "")
	if _, err := claveDesdeEntorno(); !errors.Is(err, ErrSinClave) {
		t.Fatalf("sin variables: %v", err)
	}
	t.Setenv("SUPABASE_SERVICE_ROLE_KEY", "srk")
	k3, err := claveDesdeEntorno()
	if err != nil || len(k3) != 32 {
		t.Fatalf("service role: %v %d", err, len(k3))
	}
	t.Setenv("JWT_SECRET", "jwt")
	k2, _ := claveDesdeEntorno()
	if string(k2) == string(k3) {
		t.Fatal("JWT_SECRET debe tener prioridad sobre la service role")
	}
	t.Setenv("PRESIGN_SECRET", "propia")
	k1, _ := claveDesdeEntorno()
	if string(k1) != "propia" {
		t.Fatalf("PRESIGN_SECRET se usa tal cual: %q", k1)
	}
	// La clave derivada nunca es el secreto original.
	if string(k2) == "jwt" {
		t.Fatal("la clave derivada no debe ser el secreto sin transformar")
	}
}

func TestWB_SinClave_TodasLasFuncionesFallan(t *testing.T) {
	sinClave()
	defer conClave()
	if _, _, err := Firmar("p", url.Values{}, time.Minute, t0); !errors.Is(err, ErrSinClave) {
		t.Errorf("Firmar: %v", err)
	}
	if err := Verificar("p", url.Values{}, t0); !errors.Is(err, ErrSinClave) {
		t.Errorf("Verificar: %v", err)
	}
	if _, err := EmitirToken("a", time.Minute, t0); !errors.Is(err, ErrSinClave) {
		t.Errorf("EmitirToken: %v", err)
	}
	if _, err := LeerToken("a.1.b", t0); !errors.Is(err, ErrSinClave) {
		t.Errorf("LeerToken: %v", err)
	}
}

// exp no numérico con firma correcta: camino ParseInt con error.
func TestWB_Verificar_ExpNoNumerico(t *testing.T) {
	clave, _ := ClaveFn()
	q := url.Values{"recurso": {"resumen"}, "exp": {"mañana"}}
	q.Set("sig", firma(clave, "consulta", canonico(q)))
	if err := Verificar("consulta", q, t0); !errors.Is(err, ErrFirmaInvalida) {
		t.Fatalf("exp no numérico: %v", err)
	}
}

// Valor límite: en el instante exacto de expiración todavía vale; un
// segundo después no (la condición es ahora > exp).
func TestWB_Verificar_LimiteDeExpiracion(t *testing.T) {
	q, expira, _ := Firmar("consulta", url.Values{}, 5*time.Minute, t0)
	if err := Verificar("consulta", q, expira); err != nil {
		t.Fatalf("en el segundo exacto debe valer: %v", err)
	}
	if err := Verificar("consulta", q, expira.Add(time.Second)); !errors.Is(err, ErrExpirada) {
		t.Fatalf("un segundo después debe expirar: %v", err)
	}
}

// Un exp alterado se reporta como firma inválida y NO como expirada: no hay
// forma de revivir una URL cambiando su fecha.
func TestWB_Verificar_ExpAlteradoNoRevive(t *testing.T) {
	q, _, _ := Firmar("consulta", url.Values{}, time.Second, t0)
	q.Set("exp", strconv.FormatInt(t0.Add(24*time.Hour).Unix(), 10))
	if err := Verificar("consulta", q, t0.Add(time.Hour)); !errors.Is(err, ErrFirmaInvalida) {
		t.Fatalf("esperaba firma inválida, llegó %v", err)
	}
}

func TestWB_LeerToken_Formas(t *testing.T) {
	clave, _ := ClaveFn()
	casos := map[string]string{
		"dos partes":  "a.b",
		"cuatro":      "a.1.b.c",
		"id vacío":    ".1.x",
		"exp no num":  "a.xx." + firma(clave, "soporte-token", "a.xx"),
		"firma ajena": "a.99999999999.AAAA",
	}
	for nombre, tok := range casos {
		if _, err := LeerToken(tok, t0); !errors.Is(err, ErrFirmaInvalida) {
			t.Errorf("%s: %v", nombre, err)
		}
	}
	// El token no sirve como firma de URL: propósitos distintos.
	tok, _ := EmitirToken("abc", time.Hour, t0)
	if id, err := LeerToken(tok, t0); err != nil || id != "abc" {
		t.Fatalf("token válido: %q %v", id, err)
	}
}

// Firmar no debe modificar el mapa de campos que recibe.
func TestWB_Firmar_NoMutaLaEntrada(t *testing.T) {
	campos := url.Values{"recurso": {"resumen"}, "sig": {"falsa"}}
	q, _, _ := Firmar("consulta", campos, time.Minute, t0)
	if campos.Get("exp") != "" {
		t.Fatal("Firmar mutó la entrada")
	}
	if q.Get("sig") == "falsa" {
		t.Fatal("una sig recibida debe sustituirse por la calculada")
	}
}
