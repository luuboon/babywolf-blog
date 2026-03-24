import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { Router, RouterModule } from '@angular/router';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../../environments/environment';
import { StorageService } from '../../../core/services/storage.service';

@Component({
  selector: 'app-post-editor',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterModule],
  templateUrl: './editor.component.html',
  styleUrls: ['../admin-dashboard.component.scss']
})
export class PostEditorComponent {
  title = '';
  content = '';
  coverUrl = '';
  
  isSubmitting = false;
  uploadingImage = false;
  errorMsg = '';

  private http = inject(HttpClient);
  private router = inject(Router);
  private storageService = inject(StorageService);

  onFileSelected(event: any) {
    const file: File = event.target.files[0];
    if (file) {
      this.uploadingImage = true;
      this.storageService.uploadFile(file, 'posts').subscribe({
        next: (res) => {
          this.coverUrl = res.url;
          this.uploadingImage = false;
        },
        error: (err) => {
          this.errorMsg = 'Error al subir imagen.';
          this.uploadingImage = false;
          console.error(err);
        }
      });
    }
  }

  savePost() {
    if (!this.title || !this.content) {
      this.errorMsg = 'Título y contenido son requeridos.';
      return;
    }

    this.isSubmitting = true;
    const postData = {
      title: this.title,
      content: this.content,
      cover_image: this.coverUrl,
      published: true
    };

    this.http.post(`${environment.apiUrl}/admin/posts`, postData).subscribe({
      next: () => {
        this.router.navigate(['/admin/posts']);
      },
      error: (err) => {
        this.errorMsg = 'Error al guardar el post.';
        this.isSubmitting = false;
        console.error(err);
      }
    });
  }
}
