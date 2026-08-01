import { Component, OnInit, inject, ChangeDetectorRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { SupabaseService } from '../../../core/services/supabase.service';

interface AuditEntry {
  id: string;
  email: string;
  action: string;
  ip_address: string | null;
  details: string | null;
  created_at: string;
}

const ACTION_LABELS: Record<string, string> = {
  login_success: '🟢 Inicio de sesión',
  login_failed: '🔴 Intento fallido',
  logout: '🚪 Cierre de sesión',
  password_change: '🔑 Cambio de contraseña',
  user_created: '➕ Alta de usuario',
  user_updated: '✏️ Edición de usuario',
  user_activated: '✅ Activación de cuenta',
  user_deactivated: '🚫 Desactivación de cuenta',
  role_changed: '🛡️ Cambio de rol'
};

/**
 * Practica 11-12, Parte 3: consulta del historial de accesos (bitácora).
 */
@Component({
  selector: 'app-audit-log',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './audit-log.component.html',
  styleUrls: ['../users/user-management.component.scss']
})
export class AuditLogComponent implements OnInit {
  entries: AuditEntry[] = [];
  loading = true;
  errorMsg = '';

  private sb = inject(SupabaseService);
  private cdr = inject(ChangeDetectorRef);

  ngOnInit() {
    this.load();
  }

  actionLabel(action: string): string {
    return ACTION_LABELS[action] || action;
  }

  async load() {
    this.loading = true;
    const { data, error } = await this.sb.client
      .from('audit_log')
      .select('*')
      .order('created_at', { ascending: false })
      .limit(200);

    if (error) {
      this.errorMsg = 'Error al cargar la bitácora: ' + error.message;
    } else {
      this.entries = data as AuditEntry[];
    }
    this.loading = false;
    this.cdr.markForCheck();
  }
}
