import { Component, inject, OnInit, ChangeDetectorRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { AuthService } from '../../../core/services/auth.service';
import { SupabaseService } from '../../../core/services/supabase.service';
import { StorageService } from '../../../core/services/storage.service';

@Component({
  selector: 'app-profile-page',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterModule],
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
            
            <div class="avatar-section">
              <div class="avatar-preview">
                @if (profile.avatar_url) {
                  <img [src]="profile.avatar_url" alt="Avatar de usuario" class="avatar-img" />
                } @else {
                  <div class="avatar-placeholder">👤</div>
                }
              </div>
              <div class="avatar-actions">
                <label class="btn-secondary upload-btn" [class.disabled]="uploadingAvatar">
                  {{ uploadingAvatar ? '⏳ Subiendo...' : '📷 Cambiar Foto' }}
                  <input type="file" accept="image/*" class="hidden-file-input" (change)="onAvatarSelected($event)" [disabled]="uploadingAvatar" />
                </label>
              </div>
            </div>

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
              @if (!showDisableConfirm) {
                <button class="btn-danger" (click)="showDisableConfirm = true" [disabled]="mfaLoading">
                  🗑️ Desactivar 2FA
                </button>
              } @else {
                <div class="mfa-enrollment">
                  <p class="step-label">🔐 Ingresa tu código de Google Authenticator para confirmar:</p>
                  <input type="text" [(ngModel)]="disableCode" class="field-input code-input" placeholder="000000" maxlength="6" />
                  <div class="mfa-actions">
                    <button class="btn-danger" (click)="disableMfa()" [disabled]="mfaLoading || disableCode.length < 6">
                      {{ mfaLoading ? '⏳ Verificando...' : '🗑️ Confirmar Desactivar' }}
                    </button>
                    <button class="btn-secondary" (click)="showDisableConfirm = false; disableCode = ''" [disabled]="mfaLoading">
                      Cancelar
                    </button>
                  </div>
                </div>
              }
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

          <hr class="divider" />

          <!-- SECCION MIS POSTS -->
          <h3 class="section-title">📝 Mis Posts</h3>
          <div class="my-posts-section">
            <div class="posts-header">
              <span class="posts-count">{{ myPosts.length }} publicación{{ myPosts.length !== 1 ? 'es' : '' }}</span>
              <a routerLink="/editor/new" class="btn-new-post">+ Nuevo Post</a>
            </div>

            @if (loadingPosts) {
              <div class="skeleton-block"></div>
              <div class="skeleton-block short"></div>
            } @else if (myPosts.length === 0) {
              <div class="empty-posts">
                <span class="empty-icon">📭</span>
                <p>Aún no has publicado ningún post.</p>
                <a routerLink="/editor/new" class="btn-primary">✍️ Escribir mi primer post</a>
              </div>
            } @else {
              <div class="posts-list">
                @for (post of myPosts; track post.id) {
                  <div class="post-item">
                    <div class="post-info">
                      <span class="post-title">{{ post.title }}</span>
                      <div class="post-meta">
                        <span class="post-category">{{ getCategoryEmoji(post.category) }} {{ post.category }}</span>
                        <span class="post-date">{{ post.created_at | date:'dd/MM/yyyy' }}</span>
                        <span class="post-status" [class.published]="post.published">{{ post.published ? '🟢 Publicado' : '🟡 Borrador' }}</span>
                      </div>
                    </div>
                    <div class="post-actions">
                      <a [routerLink]="['/posts', post.slug]" class="btn-view" title="Ver post">👁️</a>
                      <a [routerLink]="['/editor', post.id]" class="btn-edit" title="Editar post">✏️</a>
                    </div>
                  </div>
                }
              </div>
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

    /* ─── My Posts Section ─── */
    .my-posts-section {
      background: #fff;
      border: 2px solid #1b1b1b;
      border-radius: 8px;
      padding: 20px;
      box-shadow: 3px 3px 0 #1b1b1b;
    }
    .posts-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 16px;
    }
    .posts-count {
      font-weight: 800;
      font-size: 0.85rem;
      color: #888;
      text-transform: uppercase;
      letter-spacing: 0.04em;
    }
    .btn-new-post {
      display: inline-block;
      padding: 8px 16px;
      background: #1b1b1b;
      color: #f6f2ec;
      border: 2px solid #1b1b1b;
      border-radius: 6px;
      font-weight: 900;
      font-size: 0.85rem;
      font-family: 'Courier New', monospace;
      text-decoration: none;
      box-shadow: 3px 3px 0 #e94560;
      transition: all 0.2s;
    }
    .btn-new-post:hover {
      transform: translate(-1px, -1px);
      box-shadow: 4px 4px 0 #e94560;
    }
    .empty-posts {
      text-align: center;
      padding: 30px 20px;
    }
    .empty-icon {
      font-size: 2.5rem;
      display: block;
      margin-bottom: 12px;
    }
    .empty-posts p {
      color: #888;
      margin: 0 0 16px 0;
      font-size: 0.95rem;
    }
    .posts-list {
      display: flex;
      flex-direction: column;
      gap: 10px;
    }
    .post-item {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 14px 16px;
      background: #f6f2ec;
      border: 2px solid #ddd;
      border-radius: 8px;
      transition: border-color 0.2s, transform 0.15s;
    }
    .post-item:hover {
      border-color: #e94560;
      transform: translateX(3px);
    }
    .post-info {
      flex: 1;
      min-width: 0;
    }
    .post-title {
      display: block;
      font-weight: 800;
      font-size: 1rem;
      color: #1b1b1b;
      margin-bottom: 6px;
      white-space: nowrap;
      overflow: hidden;
      text-overflow: ellipsis;
    }
    .post-meta {
      display: flex;
      gap: 12px;
      flex-wrap: wrap;
      font-size: 0.8rem;
      color: #888;
      font-weight: 600;
    }
    .post-category {
      text-transform: capitalize;
    }
    .post-status {
      font-weight: 700;
    }
    .post-status.published {
      color: #28a745;
    }
    .post-actions {
      display: flex;
      gap: 8px;
      margin-left: 12px;
      flex-shrink: 0;
    }
    .btn-view, .btn-edit {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      width: 38px;
      height: 38px;
      border: 2px solid #1b1b1b;
      border-radius: 6px;
      text-decoration: none;
      font-size: 1.1rem;
      transition: all 0.2s;
      background: #fff;
    }
    .btn-view:hover {
      background: #e8f8ff;
      border-color: #0ea5e9;
    }
    .btn-edit:hover {
      background: #fff5f7;
      border-color: #e94560;
    }
    .avatar-section {
      display: flex;
      align-items: center;
      gap: 20px;
      margin-bottom: 30px;
      padding-bottom: 20px;
      border-bottom: 1px solid rgba(255, 255, 255, 0.05);
    }
    .avatar-preview {
      width: 100px;
      height: 100px;
      border-radius: 50%;
      overflow: hidden;
      background: #0f172a;
      border: 3px solid #e94560;
      display: flex;
      align-items: center;
      justify-content: center;
      box-shadow: 0 4px 15px rgba(233, 69, 96, 0.2);
    }
    .avatar-img {
      width: 100%;
      height: 100%;
      object-fit: cover;
    }
    .avatar-placeholder {
      font-size: 3rem;
      color: rgba(255, 255, 255, 0.2);
    }
    .hidden-file-input {
      display: none;
    }
    .upload-btn {
      cursor: pointer;
      display: inline-block;
      margin: 0;
    }
    .upload-btn.disabled {
      opacity: 0.6;
      cursor: not-allowed;
      pointer-events: none;
    }
  `]
})
export class ProfilePage implements OnInit {
  private authService = inject(AuthService);
  private supabase = inject(SupabaseService);
  private storage = inject(StorageService);

  private cdr = inject(ChangeDetectorRef);

  profile: any = null;
  loading = true;
  saving = false;
  uploadingAvatar = false;
  successMessage = '';
  errorMessage = '';

  // My Posts
  myPosts: any[] = [];
  loadingPosts = true;

  // MFA State
  isMfaEnabled = false;
  mfaLoading = false;
  factorId = '';
  qrCodeDataUri = '';
  mfaCode = '';
  mfaErrorMessage = '';
  showDisableConfirm = false;
  disableCode = '';

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

    // Load user posts (non-blocking)
    this.loadMyPosts(userId);
  }

  async loadMyPosts(userId: string) {
    this.loadingPosts = true;
    this.cdr.detectChanges();

    const { data, error } = await this.supabase.client
      .from('posts')
      .select('id, title, slug, category, published, created_at')
      .eq('author_id', userId)
      .order('created_at', { ascending: false });

    if (!error && data) {
      this.myPosts = data;
    }
    this.loadingPosts = false;
    this.cdr.detectChanges();
  }

  getCategoryEmoji(cat: string): string {
    const map: Record<string, string> = {
      gaming: '🎮', tech: '💻', opinion: '💬',
      guides: '📖', news: '📰'
    };
    return map[cat] || '📄';
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
      .update({ username: this.profile.username, avatar_url: this.profile.avatar_url })
      .eq('id', this.profile.id);

    if (error) {
      this.errorMessage = 'Error al actualizar: ' + error.message;
    } else {
      this.successMessage = '¡Perfil actualizado correctamente! ⭐';
    }
    this.saving = false;
    this.cdr.detectChanges();
  }

  onAvatarSelected(event: any) {
    const file = event.target.files[0];
    if (!file) return;

    this.uploadingAvatar = true;
    this.errorMessage = '';
    this.cdr.detectChanges();

    this.storage.uploadFile(file, 'blog_assets').subscribe({
      next: (res) => {
        this.profile.avatar_url = res.url;
        this.uploadingAvatar = false;
        // Auto guardamos
        this.updateProfile();
      },
      error: (err) => {
        this.errorMessage = 'Error al subir la imagen: ' + err.message;
        this.uploadingAvatar = false;
        this.cdr.detectChanges();
      }
    });
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
    this.mfaLoading = true;
    this.mfaErrorMessage = '';
    this.cdr.detectChanges();

    // Step 1: Verify TOTP code to elevate session to AAL2
    const { error: verifyError } = await this.authService.verifyMfaLogin(this.factorId, this.disableCode);
    if (verifyError) {
      this.mfaErrorMessage = 'Código incorrecto. Verifica e intenta de nuevo.';
      this.mfaLoading = false;
      this.cdr.detectChanges();
      return;
    }

    // Step 2: Now unenroll (requires AAL2)
    const { error } = await this.authService.unenrollMfa(this.factorId);
    if (error) {
      this.mfaErrorMessage = 'No se pudo desactivar: ' + error.message;
    } else {
      this.isMfaEnabled = false;
      this.factorId = '';
      this.mfaCode = '';
      this.disableCode = '';
      this.showDisableConfirm = false;
      this.successMessage = 'Seguridad 2FA desactivada.';
    }
    this.mfaLoading = false;
    this.cdr.detectChanges();
  }
}
