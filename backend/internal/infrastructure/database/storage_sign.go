package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// URLs preasignadas de Supabase Storage para la mesa de ayuda (Práctica 7).
//
// El bucket es privado: nadie puede leer ni escribir sin una URL firmada.
// El backend, que tiene la service role key, firma una URL de subida o de
// descarga y se la entrega a la app; la app habla directo con Storage usando
// sólo esa URL, que expira.

var storageHTTP = &http.Client{Timeout: 15 * time.Second}

func escaparRuta(ruta string) string {
	partes := strings.Split(ruta, "/")
	for i, p := range partes {
		partes[i] = url.PathEscape(p)
	}
	return strings.Join(partes, "/")
}

func storagePost(endpoint string, body any, destino any) error {
	base := os.Getenv("SUPABASE_URL")
	key := os.Getenv("SUPABASE_SERVICE_ROLE_KEY")
	if base == "" || key == "" {
		return fmt.Errorf("storage no configurado")
	}
	b, _ := json.Marshal(body)
	req, err := http.NewRequest(http.MethodPost, base+"/storage/v1"+endpoint, bytes.NewReader(b))
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("apikey", key)
	req.Header.Set("Authorization", "Bearer "+key)
	resp, err := storageHTTP.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 300 {
		msg, _ := io.ReadAll(io.LimitReader(resp.Body, 512))
		return fmt.Errorf("storage respondió %d: %s", resp.StatusCode, string(msg))
	}
	return json.NewDecoder(resp.Body).Decode(destino)
}

// FirmarSubida devuelve la URL absoluta a la que la app hace PUT del archivo.
func FirmarSubida(bucket, ruta string) (string, error) {
	var out struct {
		URL string `json:"url"`
	}
	if err := storagePost("/object/upload/sign/"+bucket+"/"+escaparRuta(ruta), map[string]any{}, &out); err != nil {
		return "", err
	}
	if out.URL == "" {
		return "", fmt.Errorf("storage no devolvió la URL de subida")
	}
	return os.Getenv("SUPABASE_URL") + "/storage/v1" + out.URL, nil
}

// FirmarDescarga devuelve una URL de lectura válida por "segundos".
func FirmarDescarga(bucket, ruta string, segundos int) (string, error) {
	var out struct {
		SignedURL string `json:"signedURL"`
	}
	if err := storagePost("/object/sign/"+bucket+"/"+escaparRuta(ruta), map[string]any{"expiresIn": segundos}, &out); err != nil {
		return "", err
	}
	if out.SignedURL == "" {
		return "", fmt.Errorf("storage no devolvió la URL de descarga")
	}
	return os.Getenv("SUPABASE_URL") + "/storage/v1" + out.SignedURL, nil
}
