package repositories

import (
	"context"
	"database/sql"
	"errors"
	"strconv"
	"strings"
	"time"

	"gorm.io/gorm"

	"babywolf-blog-backend/internal/domain/entities"
)

// GormSoporteRepository persiste la mesa de ayuda en las tablas soporte_*
// de Supabase (Postgres). Ver sql/005_soporte.sql.
type GormSoporteRepository struct {
	db *gorm.DB
}

func NewGormSoporteRepository(db *gorm.DB) *GormSoporteRepository {
	return &GormSoporteRepository{db: db}
}

var ErrInvitacionUsada = errors.New("la invitación ya fue usada o expiró")

// EsDuplicado detecta la violación de llave única de Postgres (SQLSTATE 23505).
func EsDuplicado(err error) bool {
	return err != nil && strings.Contains(err.Error(), "23505")
}

// ── Participantes ──────────────────────────────────────────────

func (r *GormSoporteRepository) CrearParticipante(ctx context.Context, p *entities.SoporteParticipante) error {
	return r.db.WithContext(ctx).Create(p).Error
}

func (r *GormSoporteRepository) Participante(ctx context.Context, id string) (*entities.SoporteParticipante, error) {
	var p entities.SoporteParticipante
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&p).Error; err != nil {
		return nil, err
	}
	return &p, nil
}

// AdminDeUsuario reutiliza el participante del admin si ya entró antes con la
// misma cuenta del blog, para no duplicarlo en cada inicio de sesión.
func (r *GormSoporteRepository) AdminDeUsuario(ctx context.Context, userID string) (*entities.SoporteParticipante, error) {
	var p entities.SoporteParticipante
	err := r.db.WithContext(ctx).Where("user_id = ? AND rol = ?", userID, "admin").First(&p).Error
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func (r *GormSoporteRepository) ActualizarParticipante(ctx context.Context, id string, cols map[string]any) error {
	return r.db.WithContext(ctx).Model(&entities.SoporteParticipante{}).Where("id = ?", id).Updates(cols).Error
}

// Latido marca actividad. nil = desconexión explícita ("Salir").
func (r *GormSoporteRepository) Latido(ctx context.Context, id string, t *time.Time) error {
	return r.db.WithContext(ctx).Model(&entities.SoporteParticipante{}).
		Where("id = ?", id).Update("ultimo_latido", t).Error
}

func (r *GormSoporteRepository) Participantes(ctx context.Context) ([]entities.SoporteParticipante, error) {
	var ps []entities.SoporteParticipante
	err := r.db.WithContext(ctx).Where("activo = ?", true).Order("created_at asc").Find(&ps).Error
	return ps, err
}

// ── Invitaciones (URL preasignada de rol) ──────────────────────

func (r *GormSoporteRepository) CrearInvitacion(ctx context.Context, inv *entities.SoporteInvitacion) error {
	return r.db.WithContext(ctx).Create(inv).Error
}

func (r *GormSoporteRepository) Invitacion(ctx context.Context, id string) (*entities.SoporteInvitacion, error) {
	var inv entities.SoporteInvitacion
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&inv).Error; err != nil {
		return nil, err
	}
	return &inv, nil
}

// Unirse crea el participante y consume la invitación en una sola
// transacción. El UPDATE condicionado (usada_por IS NULL) garantiza que dos
// dispositivos no puedan usar la misma URL al mismo tiempo.
func (r *GormSoporteRepository) Unirse(ctx context.Context, invID string, p *entities.SoporteParticipante, ahora time.Time) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(p).Error; err != nil {
			return err
		}
		res := tx.Model(&entities.SoporteInvitacion{}).
			Where("id = ? AND usada_por IS NULL AND expira_en > ?", invID, ahora).
			Update("usada_por", p.ID)
		if res.Error != nil {
			return res.Error
		}
		if res.RowsAffected != 1 {
			return ErrInvitacionUsada
		}
		return nil
	})
}

// ── Reportes ───────────────────────────────────────────────────

// MayorNumero devuelve el mayor número usado en los ID REPnnn (activos o no,
// para no reciclar identificadores de reportes eliminados).
func (r *GormSoporteRepository) MayorNumero(ctx context.Context) (int, error) {
	var ids []string
	if err := r.db.WithContext(ctx).Model(&entities.SoporteReporte{}).Pluck("id", &ids).Error; err != nil {
		return 0, err
	}
	mayor := 0
	for _, id := range ids {
		if n, err := strconv.Atoi(strings.TrimPrefix(id, "REP")); err == nil && n > mayor {
			mayor = n
		}
	}
	return mayor, nil
}

// ExisteReporte incluye los dados de baja: un ID eliminado no se reutiliza.
func (r *GormSoporteRepository) ExisteReporte(ctx context.Context, id string) (bool, error) {
	var n int64
	err := r.db.WithContext(ctx).Model(&entities.SoporteReporte{}).Where("id = ?", id).Count(&n).Error
	return n > 0, err
}

func (r *GormSoporteRepository) CrearReporte(ctx context.Context, rep *entities.SoporteReporte) error {
	return r.db.WithContext(ctx).Create(rep).Error
}

// Reporte sólo devuelve registros activos: uno eliminado se reporta como
// inexistente (TC05 y TC06).
func (r *GormSoporteRepository) Reporte(ctx context.Context, id string) (*entities.SoporteReporte, error) {
	var rep entities.SoporteReporte
	if err := r.db.WithContext(ctx).Where("id = ? AND activo = ?", id, true).First(&rep).Error; err != nil {
		return nil, err
	}
	return &rep, nil
}

func (r *GormSoporteRepository) Reportes(ctx context.Context, estado, texto string) ([]entities.SoporteReporte, error) {
	q := r.db.WithContext(ctx).Where("activo = ?", true)
	if estado != "" {
		q = q.Where("estado = ?", estado)
	}
	if t := strings.TrimSpace(texto); t != "" {
		like := "%" + strings.ToLower(t) + "%"
		q = q.Where("(LOWER(id) LIKE ? OR LOWER(equipo) LIKE ? OR LOWER(descripcion) LIKE ?)", like, like, like)
	}
	var reps []entities.SoporteReporte
	err := q.Order("updated_at desc").Limit(300).Find(&reps).Error
	return reps, err
}

func (r *GormSoporteRepository) ActualizarReporte(ctx context.Context, id string, cols map[string]any) error {
	cols["updated_at"] = time.Now().UTC()
	return r.db.WithContext(ctx).Model(&entities.SoporteReporte{}).
		Where("id = ? AND activo = ?", id, true).Updates(cols).Error
}

// AsumirDeDesconectados pasa al admin los reportes abiertos de agentes que
// perdieron la conexión.
func (r *GormSoporteRepository) AsumirDeDesconectados(ctx context.Context, adminID string, limite time.Time) (int64, error) {
	sub := r.db.Model(&entities.SoporteParticipante{}).Select("id").
		Where("rol = ? AND (ultimo_latido IS NULL OR ultimo_latido < ? OR activo = ?)", "agente", limite, false)
	res := r.db.WithContext(ctx).Model(&entities.SoporteReporte{}).
		Where("activo = ? AND estado <> ? AND agente_id IN (?)", true, "resuelto", sub).
		Updates(map[string]any{"agente_id": adminID, "estado": "asignado", "updated_at": time.Now().UTC()})
	return res.RowsAffected, res.Error
}

func (r *GormSoporteRepository) UltimoCambioReportes(ctx context.Context) (*time.Time, error) {
	var t sql.NullTime
	err := r.db.WithContext(ctx).Model(&entities.SoporteReporte{}).Select("MAX(updated_at)").Row().Scan(&t)
	if err != nil || !t.Valid {
		return nil, err
	}
	return &t.Time, nil
}

// ── Mensajes ───────────────────────────────────────────────────

func visibles(q *gorm.DB, yo string, verTodo bool) *gorm.DB {
	if verTodo {
		return q
	}
	return q.Where("(para_id IS NULL OR para_id = ? OR emisor_id = ?)", yo, yo)
}

func (r *GormSoporteRepository) CrearMensaje(ctx context.Context, m *entities.SoporteMensaje) error {
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *GormSoporteRepository) Mensajes(ctx context.Context, yo string, verTodo bool, desde int64) ([]entities.SoporteMensaje, error) {
	q := visibles(r.db.WithContext(ctx).Model(&entities.SoporteMensaje{}), yo, verTodo).Where("id > ?", desde)
	var ms []entities.SoporteMensaje
	err := q.Order("id asc").Limit(300).Find(&ms).Error
	return ms, err
}

type ConteoEmisor struct {
	EmisorID     string `json:"emisor_id"`
	EmisorNombre string `json:"emisor_nombre"`
	Total        int64  `json:"total"`
}

func (r *GormSoporteRepository) ConteoPorEmisor(ctx context.Context, yo string, verTodo bool) ([]ConteoEmisor, error) {
	var out []ConteoEmisor
	q := visibles(r.db.WithContext(ctx).Model(&entities.SoporteMensaje{}), yo, verTodo)
	err := q.Select("emisor_id, MAX(emisor_nombre) AS emisor_nombre, COUNT(*) AS total").
		Group("emisor_id").Scan(&out).Error
	return out, err
}

func (r *GormSoporteRepository) UltimoMensajeID(ctx context.Context) (int64, error) {
	var id sql.NullInt64
	err := r.db.WithContext(ctx).Model(&entities.SoporteMensaje{}).Select("MAX(id)").Row().Scan(&id)
	return id.Int64, err
}

func (r *GormSoporteRepository) BorrarMensajes(ctx context.Context) error {
	return r.db.WithContext(ctx).Where("1 = 1").Delete(&entities.SoporteMensaje{}).Error
}

// ArchivoVisible: sólo se firma la descarga de un archivo que aparece en un
// mensaje que el solicitante puede ver.
func (r *GormSoporteRepository) ArchivoVisible(ctx context.Context, ruta, yo string, verTodo bool) (bool, error) {
	var n int64
	q := visibles(r.db.WithContext(ctx).Model(&entities.SoporteMensaje{}), yo, verTodo).Where("archivo_path = ?", ruta)
	err := q.Count(&n).Error
	return n > 0, err
}
