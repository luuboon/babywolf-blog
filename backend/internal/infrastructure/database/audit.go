package database

// LogAuditAction inserta un registro en la bitácora (Practica 11-12, Parte 3)
// para acciones administrativas realizadas a través del backend Go
// (alta de usuarios, activar/desactivar, cambio de rol). ip puede ir vacío
// si no está disponible.
func LogAuditAction(userID, email, action, details, ip string) {
	if DB == nil {
		return
	}
	DB.Exec(
		`INSERT INTO public.audit_log (user_id, email, action, details, ip_address) VALUES (?, ?, ?, ?, ?)`,
		userID, email, action, details, ip,
	)
}
