import { Component, inject, OnInit, ChangeDetectorRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router, RouterModule, ActivatedRoute } from '@angular/router';
import { AuthService } from '../../../core/services/auth.service';
import { StorageService } from '../../../core/services/storage.service';
import { SupabasePostRepository } from '../../posts/infrastructure/repositories/supabase-post.repository';
import { SupabaseService } from '../../../core/services/supabase.service';

@Component({
  selector: 'app-user-editor',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterModule],
  template: `
    <div class="editor-page">
      <div class="editor-header">
        <a routerLink="/profile" class="back-link">← Volver a Mi Perfil</a>
        <h1 class="editor-title">{{ isEditMode ? '✏️ Editar Post' : '✍️ Escribir Nuevo Post' }}</h1>
        <p class="editor-subtitle">{{ isEditMode ? 'Modifica tu publicación' : 'Comparte tu contenido con la comunidad' }}</p>
      </div>

      <div *ngIf="loadingPost" class="loading-msg">⏳ Cargando post...</div>

      <div *ngIf="errorMsg" class="editor-error">
        ⚠️ {{ errorMsg }}
      </div>

      <div *ngIf="successMsg" class="editor-success">
        ✅ {{ successMsg }}
      </div>

      <div class="editor-card" *ngIf="!loadingPost">
        <div class="field">
          <label class="field-label">Título</label>
          <input type="text" [(ngModel)]="title" class="field-input" placeholder="Ej: Mi primer post en BabyWolf 🐺" />
        </div>

        <div class="field">
          <label class="field-label">Categoría</label>
          <select [(ngModel)]="category" class="field-input">
            <option value="gaming">🎮 Gaming</option>
            <option value="tech">💻 Tech</option>
            <option value="opinion">💬 Opinión</option>
            <option value="guides">📖 Guías</option>
            <option value="news">📰 Noticias</option>
          </select>
        </div>

        <div class="field">
          <label class="field-label">Imagen de Portada</label>
          <div class="upload-zone">
            <input type="file" (change)="onFileSelected($event)" accept="image/*" class="file-hidden" [disabled]="uploadingImage" id="user-cover-file" />
            <label for="user-cover-file" class="upload-label" [class.uploading]="uploadingImage">
              <span *ngIf="!uploadingImage && !coverUrl">📷 Seleccionar imagen</span>
              <span *ngIf="uploadingImage">⏳ Subiendo imagen...</span>
              <span *ngIf="coverUrl && !uploadingImage">✅ Imagen lista — Click para cambiar</span>
            </label>
          </div>
          <div *ngIf="coverUrl" class="preview">
            <img [src]="coverUrl" alt="Preview" />
          </div>
        </div>

        <div class="field">
          <label class="field-label">Contenido</label>
          <textarea [(ngModel)]="content" class="field-input field-textarea" rows="12" placeholder="Escribe tu artículo aquí..."></textarea>
        </div>

        <button class="publish-btn" (click)="savePost()" [disabled]="isSubmitting || uploadingImage">
          {{ isSubmitting ? '⏳ Guardando...' : (isEditMode ? '💾 Guardar Cambios' : '🚀 Publicar Post') }}
        </button>
      </div>
    </div>
  `,
  styles: [`
    .editor-page {
      max-width: 780px;
      margin: 0 auto;
      padding: 30px 20px;
    }
    .editor-header {
      margin-bottom: 24px;
    }
    .back-link {
      display: inline-block;
      color: #aaa;
      text-decoration: none;
      font-weight: 600;
      font-size: 0.9rem;
      margin-bottom: 12px;
      transition: color 0.2s;
    }
    .back-link:hover { color: #e94560; }
    .editor-title {
      font-size: 1.8rem;
      font-weight: 900;
      margin: 0 0 6px 0;
      color: #1b1b1b;
      font-family: 'Courier New', monospace;
    }
    .editor-subtitle {
      color: #666;
      margin: 0;
      font-size: 0.95rem;
    }
    .loading-msg {
      text-align: center;
      padding: 40px;
      font-size: 1.1rem;
      font-weight: 700;
      color: #555;
    }
    .editor-error {
      background: #fff0f0;
      border: 2px solid #e74c3c;
      border-radius: 8px;
      padding: 12px 16px;
      color: #c0392b;
      font-weight: 700;
      margin-bottom: 20px;
    }
    .editor-success {
      background: #f0fff4;
      border: 2px solid #2ecc71;
      border-radius: 8px;
      padding: 12px 16px;
      color: #27ae60;
      font-weight: 700;
      margin-bottom: 20px;
    }
    .editor-card {
      background: #f6f2ec;
      border: 3px solid #1b1b1b;
      border-radius: 10px;
      padding: 28px;
      box-shadow: 6px 6px 0 #1b1b1b;
    }
    .field {
      margin-bottom: 22px;
    }
    .field-label {
      display: block;
      font-weight: 800;
      font-size: 0.9rem;
      margin-bottom: 8px;
      color: #1b1b1b;
      text-transform: uppercase;
      letter-spacing: 0.04em;
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
      transition: border-color 0.2s, box-shadow 0.2s;
    }
    .field-input:focus {
      outline: none;
      border-color: #e94560;
      box-shadow: 0 0 0 3px rgba(233, 69, 96, 0.15);
    }
    .field-textarea {
      resize: vertical;
      min-height: 200px;
      line-height: 1.7;
    }
    select.field-input {
      cursor: pointer;
      appearance: none;
      background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' width='12' height='12' viewBox='0 0 12 12'%3E%3Cpath fill='%231b1b1b' d='M6 8L1 3h10z'/%3E%3C/svg%3E");
      background-repeat: no-repeat;
      background-position: right 14px center;
      padding-right: 36px;
    }
    .upload-zone {
      position: relative;
    }
    .file-hidden {
      position: absolute;
      width: 0;
      height: 0;
      opacity: 0;
      overflow: hidden;
    }
    .upload-label {
      display: flex;
      align-items: center;
      justify-content: center;
      padding: 20px;
      border: 2px dashed #1b1b1b;
      border-radius: 8px;
      cursor: pointer;
      font-weight: 700;
      font-size: 0.95rem;
      color: #555;
      background: #fff;
      transition: all 0.2s;
    }
    .upload-label:hover {
      border-color: #e94560;
      color: #e94560;
      background: #fff5f7;
    }
    .upload-label.uploading {
      border-color: #f39c12;
      color: #f39c12;
      cursor: wait;
    }
    .preview {
      margin-top: 12px;
      border: 2px solid #1b1b1b;
      border-radius: 8px;
      overflow: hidden;
      box-shadow: 4px 4px 0 #1b1b1b;
    }
    .preview img {
      width: 100%;
      height: 200px;
      object-fit: cover;
      display: block;
    }
    .publish-btn {
      display: block;
      width: 100%;
      padding: 14px 24px;
      background: #1b1b1b;
      color: #f6f2ec;
      border: 3px solid #1b1b1b;
      border-radius: 8px;
      font-weight: 900;
      font-size: 1.1rem;
      font-family: 'Courier New', monospace;
      cursor: pointer;
      box-shadow: 4px 4px 0 #e94560;
      transition: all 0.2s;
      margin-top: 8px;
    }
    .publish-btn:hover:not(:disabled) {
      transform: translate(-2px, -2px);
      box-shadow: 6px 6px 0 #e94560;
    }
    .publish-btn:active:not(:disabled) {
      transform: translate(2px, 2px);
      box-shadow: 2px 2px 0 #e94560;
    }
    .publish-btn:disabled {
      opacity: 0.5;
      cursor: not-allowed;
    }
  `]
})
export class UserEditorComponent implements OnInit {
  title = '';
  category = 'gaming';
  content = '';
  coverUrl = '';

  isSubmitting = false;
  uploadingImage = false;
  loadingPost = false;
  errorMsg = '';
  successMsg = '';

  isEditMode = false;
  postId = '';

  private route = inject(ActivatedRoute);
  private authService = inject(AuthService);
  private router = inject(Router);
  private storageService = inject(StorageService);
  private postRepo = inject(SupabasePostRepository);
  private supabase = inject(SupabaseService);
  private cdr = inject(ChangeDetectorRef);

  async ngOnInit() {
    const id = this.route.snapshot.paramMap.get('id');
    if (id) {
      this.isEditMode = true;
      this.postId = id;
      await this.loadPost(id);
    }
  }

  async loadPost(id: string) {
    this.loadingPost = true;
    this.cdr.detectChanges();

    const { data, error } = await this.supabase.client
      .from('posts')
      .select('*')
      .eq('id', id)
      .single();

    if (error || !data) {
      this.errorMsg = 'No se pudo cargar el post.';
      this.loadingPost = false;
      this.cdr.detectChanges();
      return;
    }

    // Verify ownership
    const userId = await this.authService.getCurrentUserId();
    if (data.author_id !== userId) {
      this.errorMsg = 'No tienes permiso para editar este post.';
      this.loadingPost = false;
      this.cdr.detectChanges();
      return;
    }

    this.title = data.title;
    this.category = data.category || 'gaming';
    this.content = data.content;
    this.coverUrl = data.cover_image_url || '';
    this.loadingPost = false;
    this.cdr.detectChanges();
  }

  onFileSelected(event: any) {
    const file: File = event.target.files[0];
    if (file) {
      this.uploadingImage = true;
      this.storageService.uploadFile(file, 'blog_assets').subscribe({
        next: (res) => {
          this.coverUrl = res.url;
          this.uploadingImage = false;
          this.cdr.detectChanges();
        },
        error: (err) => {
          this.errorMsg = 'Error al subir imagen.';
          this.uploadingImage = false;
          this.cdr.detectChanges();
          console.error(err);
        }
      });
    }
  }

  async savePost() {
    if (!this.title || !this.content) {
      this.errorMsg = 'Título y contenido son obligatorios.';
      return;
    }

    this.isSubmitting = true;
    this.errorMsg = '';
    this.cdr.detectChanges();

    const authorId = await this.authService.getCurrentUserId();
    if (!authorId) {
      this.errorMsg = 'Debes iniciar sesión para publicar.';
      this.isSubmitting = false;
      this.cdr.detectChanges();
      return;
    }

    if (this.isEditMode) {
      // UPDATE existing post
      const { error } = await this.supabase.client
        .from('posts')
        .update({
          title: this.title,
          content: this.content,
          excerpt: this.content.substring(0, 150) + '...',
          category: this.category,
          cover_image_url: this.coverUrl,
          updated_at: new Date().toISOString()
        })
        .eq('id', this.postId)
        .eq('author_id', authorId);

      if (error) {
        this.errorMsg = 'Error al actualizar: ' + error.message;
        this.isSubmitting = false;
      } else {
        this.successMsg = '¡Post actualizado exitosamente!';
        setTimeout(() => this.router.navigate(['/profile']), 1500);
      }
      this.cdr.detectChanges();
    } else {
      // CREATE new post
      const slug = this.title.toLowerCase()
        .trim()
        .replace(/[^\w\s-]/g, '')
        .replace(/[\s_-]+/g, '-')
        .replace(/^-+|-+$/g, '');

      this.postRepo.createPost({
        title: this.title,
        slug: slug,
        content: this.content,
        excerpt: this.content.substring(0, 150) + '...',
        category: this.category,
        coverImageUrl: this.coverUrl,
        authorId: authorId,
        author: 'Usuario',
        published: true,
        comments: [],
        likes: 0
      }).subscribe({
        next: () => {
          this.successMsg = '¡Post publicado exitosamente!';
          this.cdr.detectChanges();
          setTimeout(() => this.router.navigate(['/profile']), 1500);
        },
        error: (err: any) => {
          console.error(err);
          this.errorMsg = 'Error al guardar: ' + (err.message || 'Intenta de nuevo');
          this.isSubmitting = false;
          this.cdr.detectChanges();
        }
      });
    }
  }
}
