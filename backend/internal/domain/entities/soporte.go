package entities

import "time"

// Mesa de ayuda BabyWolf (Práctica 7 de Desarrollo móvil): cuatro
// dispositivos con roles administrador, agente y dos usuarios intercambian
// reportes, mensajes y archivos. Tablas creadas por sql/005_soporte.sql.

type SoporteParticipante struct {
	ID           string     `json:"id" gorm:"primaryKey"`
	Nombre       string     `json:"nombre"`
	Rol          string     `json:"rol"` // admin | agente | usuario
	Dispositivo  string     `json:"dispositivo"`
	UserID       *string    `json:"-"` // sólo el admin, que entra con su cuenta del blog
	Activo       bool       `json:"activo"`
	UltimoLatido *time.Time `json:"ultimo_latido"`
	CreatedAt    time.Time  `json:"created_at"`
}

func (SoporteParticipante) TableName() string { return "soporte_participantes" }

type SoporteInvitacion struct {
	ID        string    `json:"id" gorm:"primaryKey"`
	Rol       string    `json:"rol"`
	Nombre    string    `json:"nombre"`
	ExpiraEn  time.Time `json:"expira_en"`
	UsadaPor  *string   `json:"usada_por"`
	CreadaPor string    `json:"creada_por"`
	CreatedAt time.Time `json:"created_at"`
}

func (SoporteInvitacion) TableName() string { return "soporte_invitaciones" }

type SoporteReporte struct {
	ID            string    `json:"id" gorm:"primaryKey"`
	Equipo        string    `json:"equipo"`
	Descripcion   string    `json:"descripcion"`
	Tipo          string    `json:"tipo"`      // falla | atencion | acceso
	Prioridad     string    `json:"prioridad"` // Baja | Media | Alta
	Estado        string    `json:"estado"`    // abierto | asignado | en_proceso | resuelto
	SolicitanteID string    `json:"solicitante_id"`
	AgenteID      *string   `json:"agente_id"`
	DirectoAdmin  bool      `json:"directo_admin"`
	Activo        bool      `json:"activo"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (SoporteReporte) TableName() string { return "soporte_reportes" }

type SoporteMensaje struct {
	ID            int64     `json:"id" gorm:"primaryKey;autoIncrement"`
	ReporteID     *string   `json:"reporte_id"`
	EmisorID      string    `json:"emisor_id"`
	EmisorNombre  string    `json:"emisor_nombre"`
	EmisorRol     string    `json:"emisor_rol"`
	ParaID        *string   `json:"para_id"`
	Texto         string    `json:"texto"`
	ArchivoPath   *string   `json:"archivo_path"`
	ArchivoNombre *string   `json:"archivo_nombre"`
	CreatedAt     time.Time `json:"created_at"`
}

func (SoporteMensaje) TableName() string { return "soporte_mensajes" }
