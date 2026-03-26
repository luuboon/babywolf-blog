import { Component, inject, OnInit, ChangeDetectorRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { AuthService } from '../../../core/services/auth.service';
import { SupabaseService } from '../../../core/services/supabase.service';

@Component({
  selector: 'app-profile-page',
  standalone: true,
  imports: [CommonModule, FormsModule],
  template: `
    <div class="profile-container">
      <div class="profile-card">
        <h2 class="card-title">👾 Mi Perfil</h2>

        @if (loading) {
          <div class="skeleton-block"></div>
          <div class="skeleton-block short"></div>
          <div class="skeleton-block"></div>
        } @else if (profile) {
          <form class="profile-form" (ngSubmit)="updateProfile()">
            <div class="field">
              <label class="field-label">Email</label>
              <input type="text" [value]="profile.email" class="field-input" disabled />
              <small class="field-hint">El email no se puede cambiar.</small>
            </div>

            <div class="field">
              <label class="field-label">Nombre de Usuario</label>
              <input type="text" name="username" [(ngModel)]="profile.username" class="field-input" required />
            </div>

            <div class="field">
              <label class="field-label">Rol</label>
              <div class="role-badge" [class.admin]="profile.role === 'admin'">
                {{ profile.role | uppercase }}
              </div>
            </div>

            <button type="submit" class="btn-primary" [disabled]="saving">
              {{ saving ? '⏳ Guardando...' : '💾 Actualizar Perfil' }}
            </button>

            @if (successMessage) {
              <div class="alert success">{{ successMessage }}</div>
            }
            @if (errorMessage) {
              <div class="alert error">👾 {{ errorMessage }}</div>
            }
          </form>

          <hr class="divider" />

          <!-- SECCION SECURITY Y 2FA -->
          <h3 class="section-title">🛡️ Seguridad y 2FA</h3>
          <div class="mfa-section">
            @if (isMfaEnabled) {
              <div class="mfa-badge enabled">✅ Autenticación de Dos Pasos (2FA) Activada</div>
              <button class="btn-danger" (click)="disableMfa()" [disabled]="mfaLoading">
                {{ mfaLoading ? 'Procesando...' : '🗑️ Desactivar 2FA' }}
              </button>
            } @else {
              <div class="mfa-badge disabled">⚠️ 2FA Desactivado</div>
              
              @if (!factorId) {
                <p class="mfa-help">Habilita Google Authenticator para proteger tu cuenta.</p>
                <button class="btn-primary" (click)="startMfaEnrollment()" [disabled]="mfaLoading">
                  {{ mfaLoading ? '⏳ Cargando...' : '🔐 Configurar 2FA' }}
                </button>
              } @else {
                <div class="mfa-enrollment">
                  <p class="step-label">1. Escanea este código QR con Google Authenticator:</p>
                  <div class="qr-box">
                    <img [src]="qrCodeDataUri" alt="QR Code para Google Authenticator" class="qr-image" />
                  </div>
                  
                  <p class="step-label">2. Ingresa el código de 6 dígitos:</p>
                  <input type="text" [(ngModel)]="mfaCode" class="field-input code-input" placeholder="000000" maxlength="6" />
                  
                  <div class="mfa-actions">
                    <button class="btn-success" (click)="verifyAndEnableMfa()" [disabled]="mfaLoading || mfaCode.length < 6">
                      {{ mfaLoading ? '⏳ Verificando...' : '✅ Activar 2FA' }}
                    </button>
                    <button class="btn-secondary" (click)="cancelMfaEnrollment()" [disabled]="mfaLoading">
                      Cancelar
                    </button>
                  </div>
                </div>
              }
            }

            @if (mfaErrorMessage) {
              <div class="alert error">👾 {{ mfaErrorMessage }}</div>
            }
          </div>

        } @else {
          <div class="alert error">👾 No se pudo cargar el perfil.</div>
        }
      </div>
    </div>
  `,
  styles: [`
    .profile-container {
      max-width: 620px;
      margin: 30px auto;
      padding: 20px;
    }

    /* ─── Main Card ─── */
    .profile-card {
      padding: 28px;
      background: #f6f2ec;
      border: 3px solid #1b1b1b;
      border-radius: 10px;
      box-shadow: 6px 6px 0 #1b1b1b;
    }
    .card-title {
      font-family: 'Courier New', monospace;
      font-size: 1.8rem;
      font-weight: 900;
      margin: 0 0 20px 0;
      color: #1b1b1b;
      border-bottom: 3px dashed #1b1b1b;
      padding-bottom: 12px;
    }
    .section-title {
      font-family: 'Courier New', monospace;
      font-size: 1.3rem;
      font-weight: 900;
      margin: 0 0 16px 0;
      color: #1b1b1b;
    }
    .divider {
      border: none;
      border-top: 3px dashed #ccc;
      margin: 28px 0;
    }

    /* ─── Form ─── */
    .profile-form {
      display: flex;
      flex-direction: column;
      gap: 18px;
    }
    .field {
      display: flex;
      flex-direction: column;
    }
    .field-label {
      font-weight: 800;
      font-size: 0.85rem;
      text-transform: uppercase;
      letter-spacing: 0.04em;
      margin-bottom: 6px;
      color: #1b1b1b;
    }
    .field-input {
      width: 100%;
      padding: 12px 14px;
      border: 2px solid #1b1b1b;
      border-radius: 6px;
      font-family: 'Courier New', monospace;
      font-size: 1rem;
      background: #fff;
      box-sizing: border-box;
      color: #1b1b1b;
      transition: border-color 0.2s, box-shadow 0.2s;
    }
    .field-input:focus {
      outline: none;
      border-color: #e94560;
      box-shadow: 0 0 0 3px rgba(233, 69, 96, 0.15);
    }
    .field-input:disabled {
      background: #e8e4de;
      color: #888;
      cursor: not-allowed;
    }
    .code-input {
      text-align: center;
      font-size: 1.5rem;
      letter-spacing: 0.3em;
      max-width: 200px;
    }
    .field-hint {
      font-size: 0.8rem;
      color: #888;
      margin-top: 4px;
    }

    /* ─── Buttons ─── */
    .btn-primary, .btn-danger, .btn-success, .btn-secondary {
      display: inline-block;
      padding: 12px 20px;
      border: 2px solid #1b1b1b;
      border-radius: 6px;
      font-weight: 900;
      font-size: 1rem;
      font-family: 'Courier New', monospace;
      cursor: pointer;
      transition: all 0.2s;
      text-align: center;
    }
    .btn-primary {
      background: #1b1b1b;
      color: #f6f2ec;
      box-shadow: 4px 4px 0 #e94560;
    }
    .btn-primary:hover:not(:disabled) {
      transform: translate(-1px, -1px);
      box-shadow: 5px 5px 0 #e94560;
    }
    .btn-danger {
      background: #dc3545;
      color: white;
      box-shadow: 4px 4px 0 #1b1b1b;
    }
    .btn-danger:hover:not(:disabled) {
      background: #c0392b;
    }
    .btn-success {
      background: #28a745;
      color: white;
      box-shadow: 4px 4px 0 #1b1b1b;
    }
    .btn-success:hover:not(:disabled) {
      background: #218838;
    }
    .btn-secondary {
      background: #6c757d;
      color: white;
      box-shadow: 3px 3px 0 #1b1b1b;
    }
    button:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }
    button:active:not(:disabled) {
      transform: translate(2px, 2px);
      box-shadow: 1px 1px 0 #1b1b1b;
    }

    /* ─── Role Badge ─── */
    .role-badge {
      display: inline-block;
      padding: 6px 14px;
      border-radius: 20px;
      font-weight: 900;
      font-size: 0.85rem;
      background: #e8e4de;
      border: 2px solid #ccc;
      color: #555;
    }
    .role-badge.admin {
      background: #e94560;
      color: white;
      border-color: #c0392b;
    }

    /* ─── Alerts ─── */
    .alert {
      padding: 12px 16px;
      border-radius: 6px;
      font-weight: 700;
      font-size: 0.9rem;
      margin-top: 12px;
    }
    .alert.success {
      background: #d4edda;
      color: #155724;
      border-left: 4px solid #28a745;
    }
    .alert.error {
      background: #f8d7da;
      color: #721c24;
      border-left: 4px solid #dc3545;
    }

    /* ─── MFA Section ─── */
    .mfa-section {
      background: #fff;
      border: 2px solid #1b1b1b;
      border-radius: 8px;
      padding: 20px;
      box-shadow: 3px 3px 0 #1b1b1b;
    }
    .mfa-badge {
      font-weight: 900;
      padding: 10px 14px;
      border-radius: 6px;
      margin-bottom: 14px;
      font-size: 0.95rem;
    }
    .mfa-badge.enabled {
      background: #d4edda;
      color: #155724;
      border-left: 4px solid #28a745;
    }
    .mfa-badge.disabled {
      background: #fff3cd;
      color: #856404;
      border-left: 4px solid #f39c12;
    }
    .mfa-help {
      margin-bottom: 14px;
      font-size: 0.95rem;
      color: #555;
    }
    .mfa-enrollment {
      display: flex;
      flex-direction: column;
      gap: 12px;
    }
    .step-label {
      font-weight: 800;
      font-size: 0.9rem;
      color: #1b1b1b;
      margin: 0;
    }
    .qr-box {
      background: white;
      padding: 12px;
      border: 2px solid #1b1b1b;
      border-radius: 8px;
      align-self: center;
      width: 210px;
      height: 210px;
      display: flex;
      align-items: center;
      justify-content: center;
      box-shadow: 3px 3px 0 #1b1b1b;
    }
    .qr-image {
      width: 100%;
      height: 100%;
      object-fit: contain;
    }
    .mfa-actions {
      display: flex;
      gap: 10px;
      margin-top: 8px;
    }

    /* ─── Skeleton ─── */
    .skeleton-block {
      height: 48px;
      background: linear-gradient(90deg, #e8e4de 25%, #d4d0ca 50%, #e8e4de 75%);
      background-size: 200% 100%;
      animation: shimmer 1.5s infinite;
      border-radius: 6px;
      margin-bottom: 16px;
    }
    .skeleton-block.short {
      width: 60%;
    }
    @keyframes shimmer {
      0% { background-position: -200% 0; }
      100% { background-position: 200% 0; }
    }
  `]
})
export class ProfilePage implements OnInit {
  private authService = inject(AuthService);
  private supabase = inject(SupabaseService);

  private cdr = inject(ChangeDetectorRef);

  profile: any = null;
  loading = true;
  saving = false;
  successMessage = '';
  errorMessage = '';

  // MFA State
  isMfaEnabled = false;
  mfaLoading = false;
  factorId = '';
  qrCodeDataUri = '';
  mfaCode = '';
  mfaErrorMessage = '';

  async ngOnInit() {
    const userId = await this.authService.getCurrentUserId();
    if (!userId) {
      this.errorMessage = 'No hay sesión activa.';
      this.loading = false;
      this.cdr.detectChanges();
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
    this.cdr.detectChanges();
  }

  async checkMfaStatus() {
    const { data: factors, error } = await this.authService.getMfaFactors();
    if (error) {
      console.error('Error obteniendo MFA factors', error);
      return;
    }
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
    this.cdr.detectChanges();

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
    this.cdr.detectChanges();
  }

  async startMfaEnrollment() {
    this.mfaLoading = true;
    this.mfaErrorMessage = '';
    this.cdr.detectChanges();
    
    const { data, error } = await this.authService.enrollMfa();
    if (error) {
      this.mfaErrorMessage = 'No se pudo iniciar 2FA: ' + error.message;
      this.mfaLoading = false;
      this.cdr.detectChanges();
      return;
    }

    this.factorId = data.id;
    this.qrCodeDataUri = data.totp.qr_code;
    this.mfaLoading = false;
    this.cdr.detectChanges();
  }

  async verifyAndEnableMfa() {
    this.mfaLoading = true;
    this.mfaErrorMessage = '';
    this.cdr.detectChanges();

    const { data, error } = await this.authService.verifyMfaSetup(this.factorId, this.mfaCode);
    if (error) {
      this.mfaErrorMessage = 'Código incorrecto. Intenta de nuevo.';
      this.mfaLoading = false;
      this.cdr.detectChanges();
      return;
    }

    this.isMfaEnabled = true;
    this.mfaLoading = false;
    this.successMessage = '¡Seguridad 2FA activada con éxito!';
    this.cdr.detectChanges();
  }

  cancelMfaEnrollment() {
    this.factorId = '';
    this.mfaCode = '';
    this.qrCodeDataUri = '';
    this.mfaErrorMessage = '';
    this.cdr.detectChanges();
  }

  async disableMfa() {
    if (!confirm('¿Estás seguro que deseas desactivar la capa de seguridad 2FA?')) return;
    this.mfaLoading = true;
    this.cdr.detectChanges();
    
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
    this.cdr.detectChanges();
  }
}
