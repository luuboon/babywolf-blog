// Package security firma URLs preasignadas y tokens de sesión con HMAC-SHA256.
//
// Una URL preasignada lleva en su query los parámetros que autoriza, la fecha
// de expiración (exp) y una firma (sig). Si alguien modifica cualquier
// parámetro, la firma deja de coincidir y el backend rechaza la petición: eso
// es lo que comprueba la "Prueba 3 — Modificación de la URL" de la práctica 3.
//
// No se usa ninguna dependencia externa: sólo crypto/hmac de la biblioteca
// estándar.
package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

var (
	ErrFirmaInvalida = errors.New("firma inválida")
	ErrExpirada      = errors.New("url expirada")
	ErrSinClave      = errors.New("clave de firma no configurada")
)

// ClaveFn devuelve la clave HMAC. Es una variable para poder sustituirla en
// las pruebas sin depender de variables de entorno.
var ClaveFn = claveDesdeEntorno

// claveDesdeEntorno usa PRESIGN_SECRET si existe. Si no, deriva una clave
// propia a partir de un secreto que el backend ya tiene en Railway, para no
// obligar a configurar nada nuevo. Derivar (HMAC con una etiqueta) evita
// reutilizar el secreto original tal cual.
func claveDesdeEntorno() ([]byte, error) {
	if s := os.Getenv("PRESIGN_SECRET"); s != "" {
		return []byte(s), nil
	}
	for _, nombre := range []string{"JWT_SECRET", "SUPABASE_SERVICE_ROLE_KEY"} {
		if s := os.Getenv(nombre); s != "" {
			m := hmac.New(sha256.New, []byte(s))
			m.Write([]byte("babywolf-presign-v1"))
			return m.Sum(nil), nil
		}
	}
	return nil, ErrSinClave
}

func firma(clave []byte, proposito, contenido string) string {
	m := hmac.New(sha256.New, clave)
	m.Write([]byte(proposito))
	m.Write([]byte{'\n'})
	m.Write([]byte(contenido))
	return base64.RawURLEncoding.EncodeToString(m.Sum(nil))
}

// canonico serializa los parámetros sin "sig", ordenados por clave
// (url.Values.Encode ya ordena), para que la firma no dependa del orden.
func canonico(q url.Values) string {
	c := url.Values{}
	for k, v := range q {
		if k != "sig" {
			c[k] = v
		}
	}
	return c.Encode()
}

// Firmar agrega exp y sig a los campos y devuelve la query lista para la URL.
// El propósito separa dominios: una firma de "consulta" no sirve como
// invitación de "soporte" aunque los parámetros coincidan.
func Firmar(proposito string, campos url.Values, ttl time.Duration, ahora time.Time) (url.Values, time.Time, error) {
	clave, err := ClaveFn()
	if err != nil {
		return nil, time.Time{}, err
	}
	q := url.Values{}
	for k, v := range campos {
		q[k] = append([]string(nil), v...)
	}
	expira := ahora.Add(ttl).UTC().Truncate(time.Second)
	q.Set("exp", strconv.FormatInt(expira.Unix(), 10))
	q.Del("sig")
	q.Set("sig", firma(clave, proposito, canonico(q)))
	return q, expira, nil
}

// Verificar comprueba firma y expiración. Primero la firma: así un exp
// alterado se reporta como URL modificada y no como "expirada".
func Verificar(proposito string, q url.Values, ahora time.Time) error {
	clave, err := ClaveFn()
	if err != nil {
		return err
	}
	sig := q.Get("sig")
	if sig == "" {
		return ErrFirmaInvalida
	}
	esperada := firma(clave, proposito, canonico(q))
	if !hmac.Equal([]byte(sig), []byte(esperada)) {
		return ErrFirmaInvalida
	}
	exp, err := strconv.ParseInt(q.Get("exp"), 10, 64)
	if err != nil {
		return ErrFirmaInvalida
	}
	if ahora.Unix() > exp {
		return ErrExpirada
	}
	return nil
}

// EmitirToken crea un token de sesión "id.exp.firma" para un participante de
// la mesa de ayuda. No es un JWT: sólo lo emite y lo lee este backend.
func EmitirToken(id string, ttl time.Duration, ahora time.Time) (string, error) {
	clave, err := ClaveFn()
	if err != nil {
		return "", err
	}
	exp := strconv.FormatInt(ahora.Add(ttl).Unix(), 10)
	cuerpo := id + "." + exp
	return cuerpo + "." + firma(clave, "soporte-token", cuerpo), nil
}

// LeerToken valida el token y devuelve el id del participante.
func LeerToken(token string, ahora time.Time) (string, error) {
	clave, err := ClaveFn()
	if err != nil {
		return "", err
	}
	partes := strings.Split(token, ".")
	if len(partes) != 3 || partes[0] == "" {
		return "", ErrFirmaInvalida
	}
	cuerpo := partes[0] + "." + partes[1]
	if !hmac.Equal([]byte(partes[2]), []byte(firma(clave, "soporte-token", cuerpo))) {
		return "", ErrFirmaInvalida
	}
	exp, err := strconv.ParseInt(partes[1], 10, 64)
	if err != nil {
		return "", ErrFirmaInvalida
	}
	if ahora.Unix() > exp {
		return "", ErrExpirada
	}
	return partes[0], nil
}
