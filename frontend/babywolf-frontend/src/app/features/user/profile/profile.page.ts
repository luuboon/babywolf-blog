import { Component, inject, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { DomSanitizer, SafeHtml } from '@angular/platform-browser';
import { AuthService } from '../../../core/services/auth.service';
import { SupabaseService } from '../../../core/services/supabase.service';

@Component({
  selector: 'app-profile-page',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <div class="profile-container">
      <div class="neo-card profile-card">
        <h2 class="card-title">👾 Mi Perfil</h2>

        @if (loading) {
          <p class="loading-state">Cargando datos...</p>
        } @else if (profile) {
          <form class="neo-form" (ngSubmit)="updateProfile()">
            <div class="form-group">
              <label for="email">Email</label>
              <input type="text" id="email" [value]="profile.email" class="neo-input" disabled />
              <small class="hint">El email no se puede cambiar.</small>
            </div>

            <div class="form-group">
              <label for="username">Nombre de Usuario</label>
              <input type="text" id="username" name="username" [(ngModel)]="profile.username" class="neo-input" required />
            </div>

            <div class="form-group">
              <label>Rol</label>
              <div class="role-badge" [class.admin]="profile.role === 'admin'">
                {{ profile.role | uppercase }}
              </div>
            </div>

            <div class="form-actions">
              <button type="submit" class="neo-btn" [disabled]="saving">
                {{ saving ? 'Guardando...' : 'Actualizar Perfil' }}
              </button>
            </div>

            @if (successMessage) {
              <div class="success-msg">{{ successMessage }}</div>
            }
            @if (errorMessage) {
              <div class="error-msg">👾 {{ errorMessage }}</div>
            }
          </form>

          <hr class="neo-divider" />

          <!-- SECCION SECURITY Y 2FA -->
          <h3 class="card-subtitle">🛡️ Seguridad y 2FA</h3>
          <div class="mfa-section">
            @if (isMfaEnabled) {
              <div class="mfa-status success">✅ Autenticación de Dos Pasos (2FA) Activada</div>
              <button class="neo-btn btn-danger" (click)="disableMfa()" [disabled]="mfaLoading">
                {{ mfaLoading ? 'Procesando...' : 'Desactivar 2FA' }}
              </button>
            } @else {
              <div class="mfa-status warning">⚠️ 2FA Desactivado</div>
              
              @if (!factorId) {
                <p class="mfa-help">Habilita Google Authenticator para proteger tu cuenta.</p>
                <button class="neo-btn" (click)="startMfaEnrollment()" [disabled]="mfaLoading">
                  {{ mfaLoading ? 'Cargando...' : 'Configurar 2FA' }}
                </button>
              } @else {
                <div class="mfa-enrollment">
                  <p>1. Escanea este código QR con Google Authenticator o Microsoft Authenticator:</p>
                  <div class="qr-box" [innerHTML]="qrCodeHtml"></div>
                  
                  <p>2. Ingresa el código de 6 dígitos que aparece en tu app:</p>
                  <div class="form-group">
                    <input type="text" [(ngModel)]="mfaCode" class="neo-input text-center" placeholder="000000" maxlength="6" />
                  </div>
                  
                  <div class="mfa-actions">
                    <button class="neo-btn btn-success" (click)="verifyAndEnableMfa()" [disabled]="mfaLoading || mfaCode.length < 6">
                      {{ mfaLoading ? 'Verificando...' : 'Activar 2FA' }}
                    </button>
                    <button class="neo-btn btn-secondary" (click)="cancelMfaEnrollment()" [disabled]="mfaLoading">
                      Cancelar
                    </button>
                  </div>
                </div>
              }
            }

            @if (mfaErrorMessage) {
              <div class="error-msg">👾 {{ mfaErrorMessage }}</div>
            }
          </div>

        } @else {
          <p class="error-msg">👾 No se pudo cargar el perfil.</p>
        }
      </div>
    </div>
  `,
  styles: [`
    .profile-container {
      max-width: 600px;
      margin: 40px auto;
      padding: 20px;
    }
    .profile-card {
      padding: 30px;
      background: var(--surface-card);
      border: 2px solid var(--border-color);
      box-shadow: 4px 4px 0 var(--shadow-color);
      border-radius: 8px;
    }
    .card-title {
      font-family: var(--font-brand);
      font-size: 2rem;
      margin-top: 0;
      margin-bottom: 20px;
      color: var(--primary-color);
      border-bottom: 2px dashed var(--border-color);
      padding-bottom: 10px;
    }
    .card-subtitle {
      font-family: var(--font-brand);
      font-size: 1.5rem;
      margin-top: 0;
      margin-bottom: 15px;
      color: var(--text-color);
    }
    .neo-divider {
      border: none;
      border-top: 2px dashed var(--border-color);
      margin: 30px 0;
    }
    .neo-form {
      display: flex;
      flex-direction: column;
      gap: 20px;
    }
    .form-group label {
      display: block;
      font-weight: bold;
      margin-bottom: 8px;
    }
    .neo-input {
      width: 100%;
      padding: 12px;
      border: 2px solid var(--border-color);
      border-radius: 4px;
      font-family: var(--font-base);
      font-size: 1rem;
      box-sizing: border-box;
      background: var(--bg-color);
    }
    .neo-input.text-center {
      text-align: center;
      font-size: 1.5rem;
      letter-spacing: 0.2em;
    }
    .neo-input:disabled {
      background: #e0e0e0;
      color: #777;
      cursor: not-allowed;
    }
    .hint {
      font-size: 0.85rem;
      color: #666;
      display: block;
      margin-top: 4px;
    }
    .neo-btn {
      background: var(--primary-color);
      color: white;
      border: 2px solid var(--border-color);
      box-shadow: 2px 2px 0 var(--border-color);
      padding: 12px 20px;
      font-weight: bold;
      font-size: 1.1rem;
      border-radius: 4px;
      cursor: pointer;
      width: 100%;
      text-align: center;
      transition: all 0.2s ease;
    }
    .neo-btn:active {
      transform: translate(2px, 2px);
      box-shadow: 0 0 0 var(--border-color);
    }
    .neo-btn:disabled {
      background: #ccc;
      cursor: not-allowed;
      transform: none;
      box-shadow: 2px 2px 0 var(--border-color);
    }
    .btn-danger {
      background: #dc3545;
    }
    .btn-success {
      background: #28a745;
    }
    .btn-secondary {
      background: #6c757d;
    }
    .role-badge {
      display: inline-block;
      padding: 6px 12px;
      border-radius: 20px;
      font-weight: bold;
      font-size: 0.9rem;
      background: #eee;
      border: 2px solid #ccc;
    }
    .role-badge.admin {
      background: var(--accent-color);
      color: white;
      border-color: #c0392b;
    }
    .success-msg {
      background: #d4edda;
      color: #155724;
      padding: 10px;
      border-radius: 4px;
      border-left: 4px solid #28a745;
      margin-top: 15px;
    }
    .error-msg {
      background: #f8d7da;
      color: #721c24;
      padding: 10px;
      border-radius: 4px;
      border-left: 4px solid #dc3545;
      margin-top: 15px;
      font-weight: bold;
    }
    .mfa-section {
      background: #f9f9f9;
      border: 2px solid var(--border-color);
      padding: 20px;
      border-radius: 8px;
    }
    .mfa-status {
      font-weight: bold;
      margin-bottom: 15px;
      padding: 10px;
      border-radius: 4px;
    }
    .mfa-status.success {
      background: #d4edda;
      color: #155724;
      border-left: 4px solid #28a745;
    }
    .mfa-status.warning {
      background: #fff3cd;
      color: #856404;
      border-left: 4px solid #ffeeba;
    }
    .mfa-help {
      margin-bottom: 15px;
      font-size: 0.95rem;
    }
    .mfa-enrollment {
      display: flex;
      flex-direction: column;
      gap: 15px;
    }
    .qr-box {
      background: white;
      padding: 10px;
      border: 2px solid var(--border-color);
      border-radius: 8px;
      align-self: center;
      width: 200px;
      height: 200px;
      display: flex;
      align-items: center;
      justify-content: center;
    }
    .mfa-actions {
      display: flex;
      gap: 10px;
      margin-top: 10px;
    }
    .loading-state {
      text-align: center;
      font-weight: bold;
      padding: 30px;
      font-size: 1.2rem;
    }
  `]
})
export class ProfilePage implements OnInit {
  private authService = inject(AuthService);
  private supabase = inject(SupabaseService);
  private sanitizer = inject(DomSanitizer);

  profile: any = null;
  loading = true;
  saving = false;
  successMessage = '';
  errorMessage = '';

  // MFA State
  isMfaEnabled = false;
  mfaLoading = false;
  factorId = '';
  qrCodeHtml: SafeHtml = '';
  mfaCode = '';
  mfaErrorMessage = '';

  async ngOnInit() {
    const userId = await this.authService.getCurrentUserId();
    if (!userId) {
      this.errorMessage = 'No hay sesión activa.';
      this.loading = false;
      return;
    }

    // Load Profile
    const { data, error } = await this.supabase.client
      .from('users')
      .select('*')
      .eq('id', userId)
      .single();

    if (error) {
      this.errorMessage = 'Error al cargar el perfil: ' + error.message;
    } else {
      this.profile = data;
    }

    // Load MFA Status
    await this.checkMfaStatus();

    this.loading = false;
  }

  async checkMfaStatus() {
    const { data: factors, error } = await this.authService.getMfaFactors();
    if (error) {
      console.error('Error obteniendo MFA factors', error);
      return;
    }
    // Verificamos si hay al menos un factor de estado 'verified'
    const verifiedFactor = factors?.all.find((f: any) => f.status === 'verified' && f.factor_type === 'totp');
    if (verifiedFactor) {
      this.isMfaEnabled = true;
      this.factorId = verifiedFactor.id;
    }
  }

  async updateProfile() {
    if (!this.profile) return;
    this.saving = true;
    this.successMessage = '';
    this.errorMessage = '';

    const { error } = await this.supabase.client
      .from('users')
      .update({ username: this.profile.username })
      .eq('id', this.profile.id);

    if (error) {
      this.errorMessage = 'Error al actualizar: ' + error.message;
    } else {
      this.successMessage = '¡Perfil actualizado correctamente! ⭐';
    }
    this.saving = false;
  }

  async startMfaEnrollment() {
    this.mfaLoading = true;
    this.mfaErrorMessage = '';
    
    // Enroll via Supabase
    const { data, error } = await this.authService.enrollMfa();
    if (error) {
      this.mfaErrorMessage = 'No se pudo iniciar 2FA: ' + error.message;
      this.mfaLoading = false;
      return;
    }

    this.factorId = data.id;
    // Sanitizamos el SVG retornado
    this.qrCodeHtml = this.sanitizer.bypassSecurityTrustHtml(data.totp.qr_code);
    this.mfaLoading = false;
  }

  async verifyAndEnableMfa() {
    this.mfaLoading = true;
    this.mfaErrorMessage = '';

    const { data, error } = await this.authService.verifyMfaSetup(this.factorId, this.mfaCode);
    if (error) {
      this.mfaErrorMessage = 'Código incorrecto. Intenta de nuevo.';
      this.mfaLoading = false;
      return;
    }

    // Activación exitosa
    this.isMfaEnabled = true;
    this.mfaLoading = false;
    this.successMessage = '¡Seguridad 2FA activada con éxito!';
  }

  cancelMfaEnrollment() {
    // Si queremos ser meticulosos podríamos hacer unenroll parcial, pero simplemente limpiamos vista.
    this.factorId = '';
    this.mfaCode = '';
    this.qrCodeHtml = '';
    this.mfaErrorMessage = '';
  }

  async disableMfa() {
    if (!confirm('¿Estás seguro que deseas desactivar la capa de seguridad 2FA?')) return;
    this.mfaLoading = true;
    
    const { error } = await this.authService.unenrollMfa(this.factorId);
    if (error) {
      this.mfaErrorMessage = 'No se pudo desactivar: ' + error.message;
    } else {
      this.isMfaEnabled = false;
      this.factorId = '';
      this.mfaCode = '';
      this.successMessage = 'Seguridad 2FA desactivada.';
    }
    this.mfaLoading = false;
  }
}
