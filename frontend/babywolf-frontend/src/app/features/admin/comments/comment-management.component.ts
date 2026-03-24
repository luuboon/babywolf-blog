import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../../environments/environment';

interface Comment {
  id: string;
  content: string;
  author: { email: string };
  post: { title: string };
  created_at: string;
}

@Component({
  selector: 'app-comment-management',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './comment-management.component.html',
  styleUrls: ['../admin-dashboard.component.scss']
})
export class CommentManagementComponent implements OnInit {
  comments: Comment[] = [];
  loading = true;
  errorMsg = '';

  constructor(private http: HttpClient) {}

  ngOnInit() {
    this.loadComments();
  }

  loadComments() {
    this.loading = true;
    this.http.get<Comment[]>(`${environment.apiUrl}/admin/comments`).subscribe({
      next: (data) => {
        this.comments = data;
        this.loading = false;
      },
      error: (err) => {
        this.errorMsg = 'Error al cargar los comentarios.';
        this.loading = false;
        console.error(err);
      }
    });
  }

  deleteComment(id: string) {
    if (!confirm('¿Seguro de que deseas eliminar este comentario? Acción irreversible.')) return;
    
    this.http.delete(`${environment.apiUrl}/admin/comments/${id}`).subscribe({
      next: () => {
        this.comments = this.comments.filter(c => c.id !== id);
      },
      error: (err) => {
        alert('Error al eliminar comentario');
        console.error(err);
      }
    });
  }
}
