import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../../environments/environment';

interface Post {
  id: string;
  title: string;
  author: { email: string };
  created_at: string;
  published: boolean;
}

@Component({
  selector: 'app-post-management',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './post-management.component.html',
  styleUrls: ['../admin-dashboard.component.scss']
})
export class PostManagementComponent implements OnInit {
  posts: Post[] = [];
  loading = true;
  errorMsg = '';

  constructor(private http: HttpClient) {}

  ngOnInit() {
    this.loadPosts();
  }

  loadPosts() {
    this.loading = true;
    this.http.get<Post[]>(`${environment.apiUrl}/admin/posts`).subscribe({
      next: (data) => {
        this.posts = data;
        this.loading = false;
      },
      error: (err) => {
        this.errorMsg = 'Error al cargar posts administrables.';
        this.loading = false;
        console.error(err);
      }
    });
  }

  deletePost(id: string) {
    if (!confirm('¿Seguro de que deseas eliminar este post permanentemente?')) return;
    
    this.http.delete(`${environment.apiUrl}/admin/posts/${id}`).subscribe({
      next: () => {
        this.posts = this.posts.filter(p => p.id !== id);
      },
      error: (err) => {
        alert('Error al eliminar post');
        console.error(err);
      }
    });
  }

  togglePublish(post: Post) {
    this.http.put(`${environment.apiUrl}/admin/posts/${post.id}/publish`, { published: !post.published }).subscribe({
      next: () => {
        post.published = !post.published;
      },
      error: (err) => {
        alert('Error al actualizar estado de publicación');
        console.error(err);
      }
    });
  }
}
