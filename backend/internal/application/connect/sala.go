// Package connect implementa el canal de mensajes de MobileConnect
// (Práctica 6): dos instancias de la app, "Dispositivo A" y "Dispositivo B",
// se escriben a través de una sala identificada por un código.
//
// Vive en memoria a propósito: es un canal efímero de demostración y no
// necesita tablas nuevas en Supabase. Si Railway reinicia el servicio, el
// historial se pierde (el botón "Limpiar historial" hace lo mismo a mano).
package connect

import (
	"errors"
	"regexp"
	"strings"
	"sync"
	"time"
	"unicode/utf8"
)

const (
	MaxTexto       = 280
	MaxDispositivo = 30
	maxPorSala     = 200
	maxSalas       = 200
	vidaSala       = 24 * time.Hour
)

var (
	ErrSala        = errors.New("el código de sala debe tener de 3 a 24 letras, números o guiones")
	ErrEmisor      = errors.New("el emisor debe ser el dispositivo A o el B")
	ErrVacio       = errors.New("el mensaje no puede estar vacío")
	ErrLargo       = errors.New("el mensaje supera los 280 caracteres")
	ErrDispositivo = errors.New("el nombre del dispositivo supera los 30 caracteres")
)

var reSala = regexp.MustCompile(`^[A-Z0-9-]{3,24}$`)

type Mensaje struct {
	ID          int64     `json:"id"`
	Emisor      string    `json:"emisor"`
	Dispositivo string    `json:"dispositivo"`
	Texto       string    `json:"texto"`
	Hora        time.Time `json:"hora"`
}

type sala struct {
	mensajes  []Mensaje
	conteo    map[string]int
	actividad time.Time
}

type Store struct {
	mu    sync.Mutex
	salas map[string]*sala
	seq   int64
	Ahora func() time.Time
}

func NewStore() *Store {
	return &Store{salas: map[string]*sala{}, Ahora: time.Now}
}

// NormalizarSala acepta minúsculas y espacios alrededor.
func NormalizarSala(codigo string) (string, error) {
	c := strings.ToUpper(strings.TrimSpace(codigo))
	if !reSala.MatchString(c) {
		return "", ErrSala
	}
	return c, nil
}

func (s *Store) purgar(ahora time.Time) {
	for k, v := range s.salas {
		if ahora.Sub(v.actividad) > vidaSala {
			delete(s.salas, k)
		}
	}
	// Tope de salas: se descarta la menos activa.
	for len(s.salas) > maxSalas {
		var vieja string
		var t time.Time
		for k, v := range s.salas {
			if vieja == "" || v.actividad.Before(t) {
				vieja, t = k, v.actividad
			}
		}
		delete(s.salas, vieja)
	}
}

// Enviar valida y guarda un mensaje. Las reglas de la práctica viven aquí y
// no sólo en la app: un cliente modificado tampoco puede mandar vacíos.
func (s *Store) Enviar(codigo, emisor, dispositivo, texto string) (Mensaje, error) {
	c, err := NormalizarSala(codigo)
	if err != nil {
		return Mensaje{}, err
	}
	emisor = strings.ToUpper(strings.TrimSpace(emisor))
	if emisor != "A" && emisor != "B" {
		return Mensaje{}, ErrEmisor
	}
	texto = strings.TrimSpace(texto)
	if texto == "" {
		return Mensaje{}, ErrVacio
	}
	if utf8.RuneCountInString(texto) > MaxTexto {
		return Mensaje{}, ErrLargo
	}
	dispositivo = strings.TrimSpace(dispositivo)
	if utf8.RuneCountInString(dispositivo) > MaxDispositivo {
		return Mensaje{}, ErrDispositivo
	}
	if dispositivo == "" {
		dispositivo = "Dispositivo " + emisor
	}

	s.mu.Lock()
	defer s.mu.Unlock()
	ahora := s.Ahora().UTC()
	s.purgar(ahora)

	sl := s.salas[c]
	if sl == nil {
		sl = &sala{conteo: map[string]int{}}
		s.salas[c] = sl
	}
	s.seq++
	m := Mensaje{ID: s.seq, Emisor: emisor, Dispositivo: dispositivo, Texto: texto, Hora: ahora}
	sl.mensajes = append(sl.mensajes, m)
	if len(sl.mensajes) > maxPorSala {
		sl.mensajes = sl.mensajes[len(sl.mensajes)-maxPorSala:]
	}
	sl.conteo[emisor]++
	sl.actividad = ahora
	return m, nil
}

// Listar devuelve los mensajes con id mayor que "desde" y el contador de
// enviados por dispositivo (el contador cuenta todo, no sólo lo nuevo).
func (s *Store) Listar(codigo string, desde int64) ([]Mensaje, map[string]int, error) {
	c, err := NormalizarSala(codigo)
	if err != nil {
		return nil, nil, err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	out := []Mensaje{}
	conteo := map[string]int{"A": 0, "B": 0}
	sl := s.salas[c]
	if sl == nil {
		return out, conteo, nil
	}
	for _, m := range sl.mensajes {
		if m.ID > desde {
			out = append(out, m)
		}
	}
	for k, v := range sl.conteo {
		conteo[k] = v
	}
	return out, conteo, nil
}

// Limpiar borra el historial y los contadores de la sala.
func (s *Store) Limpiar(codigo string) error {
	c, err := NormalizarSala(codigo)
	if err != nil {
		return err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.salas, c)
	return nil
}
