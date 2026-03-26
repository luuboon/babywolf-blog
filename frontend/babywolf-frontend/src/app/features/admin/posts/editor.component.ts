import { Component, inject, ChangeDetectorRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { AuthService } from '../../../core/services/auth.service';
import { StorageService } from '../../../core/services/storage.service';
import { SupabasePostRepository } from '../../posts/infrastructure/repositories/supabase-post.repository';

@Component({
  selector: 'app-post-editor',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterModule],
  templateUrl: './editor.component.html',
  styleUrls: ['../admin-dashboard.component.scss']
})
export class PostEditorComponent {
  title = '';
  category = 'gaming';
  content = '';
  coverUrl = '';
  
  isSubmitting = false;
  uploadingImage = false;
  errorMsg = '';

  private authService = inject(AuthService);
  private router = inject(Router);
  private storageService = inject(StorageService);
  private postRepo = inject(SupabasePostRepository);
  private cdr = inject(ChangeDetectorRef);

  onFileSelected(event: any) {
    const file: File = event.target.files[0];
    if (file) {
      this.uploadingImage = true;
      this.storageService.uploadFile(file, 'blog_assets').subscribe({
        next: (res) => {
          this.coverUrl = res.url;
          this.uploadingImage = false;
          this.cdr.markForCheck();
        },
        error: (err) => {
          this.errorMsg = 'Error al subir imagen. Asegúrate que el bucket "posts" sea público.';
          this.uploadingImage = false;
          this.cdr.markForCheck();
          console.error(err);
        }
      });
    }
  }

  async savePost() {
    if (!this.title || !this.content) {
      this.errorMsg = 'Título y contenido son requeridos.';
      return;
    }

    this.isSubmitting = true;
    const authorId = await this.authService.getCurrentUserId();
    
    if (!authorId) {
      this.errorMsg = 'Debes estar autenticado para publicar.';
      this.isSubmitting = false;
      this.cdr.markForCheck();
      return;
    }

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
      author: 'Admin',
      published: true,
      comments: [],
      likes: 0
    }).subscribe({
      next: () => {
        this.router.navigate(['/admin/posts']);
      },
      error: (err: any) => {
        console.error(err);
        this.errorMsg = 'Error al guardar el post: ' + (err.message || 'Error desconocido');
        this.isSubmitting = false;
      }
    });
  }
}
