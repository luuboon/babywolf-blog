package middleware

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"math/big"
	"net/http"
	"os"
	"sync"
	"time"
)

// Soporte para las JWT Signing Keys asimétricas (ECC P-256 / ES256) de Supabase.
// Desde 2025 los proyectos nuevos firman los access_token con una llave asimétrica
// y publican la parte pública en el endpoint JWKS. Aquí la descargamos, la
// convertimos a *ecdsa.PublicKey y la cacheamos.

type jwk struct {
	Kid string `json:"kid"`
	Kty string `json:"kty"`
	Crv string `json:"crv"`
	X   string `json:"x"`
	Y   string `json:"y"`
}

type jwksResponse struct {
	Keys []jwk `json:"keys"`
}

var (
	jwksCache   = map[string]*ecdsa.PublicKey{}
	jwksMu      sync.RWMutex
	jwksFetched time.Time
)

// getECDSAPublicKey devuelve la llave pública para el kid dado, refrescando el
// JWKS si no está en caché o si la caché tiene más de una hora.
// ponytail: caché en memoria de proceso; suficiente para un solo servicio.
func getECDSAPublicKey(kid string) (*ecdsa.PublicKey, error) {
	jwksMu.RLock()
	key, ok := jwksCache[kid]
	fresh := time.Since(jwksFetched) < time.Hour
	jwksMu.RUnlock()
	if ok && fresh {
		return key, nil
	}

	if err := refreshJWKS(); err != nil {
		return nil, err
	}

	jwksMu.RLock()
	defer jwksMu.RUnlock()
	if key, ok := jwksCache[kid]; ok {
		return key, nil
	}
	return nil, fmt.Errorf("kid %q not found in JWKS", kid)
}

func refreshJWKS() error {
	url := os.Getenv("SUPABASE_URL") + "/auth/v1/.well-known/jwks.json"
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var out jwksResponse
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return err
	}

	newCache := map[string]*ecdsa.PublicKey{}
	for _, k := range out.Keys {
		if k.Kty != "EC" || k.Crv != "P-256" {
			continue
		}
		xb, err := base64.RawURLEncoding.DecodeString(k.X)
		if err != nil {
			continue
		}
		yb, err := base64.RawURLEncoding.DecodeString(k.Y)
		if err != nil {
			continue
		}
		newCache[k.Kid] = &ecdsa.PublicKey{
			Curve: elliptic.P256(),
			X:     new(big.Int).SetBytes(xb),
			Y:     new(big.Int).SetBytes(yb),
		}
	}

	jwksMu.Lock()
	jwksCache = newCache
	jwksFetched = time.Now()
	jwksMu.Unlock()
	return nil
}
