import { Component, OnInit, inject } from '@angular/core';
import { CommonModule, DatePipe } from '@angular/common';
import { ActivatedRoute, RouterModule } from '@angular/router';
import { Post } from '../../../domain/models/post.model';
import { POST_REPOSITORY } from '../../../domain/repositories/post.repository.token';

@Component({
  selector: 'app-post-detail',
  standalone: true,
  imports: [CommonModule, RouterModule, DatePipe],
  templateUrl: './post-detail.html',
  styleUrls: ['./post-detail.scss']
})
export class PostDetailPage implements OnInit {
  post: Post | null = null;
  isLoading = true;
  error = '';

  private route = inject(ActivatedRoute);
  private postRepo = inject(POST_REPOSITORY);

  ngOnInit(): void {
    const slug = this.route.snapshot.paramMap.get('slug') || '';
    this.postRepo.getPostBySlug(slug).subscribe({
      next: (post) => {
        this.post = post || null;
        this.isLoading = false;
        if (!post) {
          this.error = 'Post no encontrado.';
        }
      },
      error: () => {
        this.error = 'Error al cargar el post.';
        this.isLoading = false;
      }
    });
  }
}
