import { Component, OnInit, inject, ChangeDetectorRef } from '@angular/core';
import { CommonModule, DatePipe } from '@angular/common';
import { ActivatedRoute, RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { Post } from '../../../domain/models/post.model';
import { POST_REPOSITORY } from '../../../domain/repositories/post.repository.token';
import { AuthService } from '../../../../../core/services/auth.service';

@Component({
  selector: 'app-post-detail',
  standalone: true,
  imports: [CommonModule, RouterModule, DatePipe, FormsModule],
  templateUrl: './post-detail.html',
  styleUrls: ['./post-detail.scss']
})
export class PostDetailPage implements OnInit {
  post: Post | null = null;
  isLoading = true;
  error = '';

  newCommentText = '';
  isSubmittingComment = false;
  commentError = '';

  private route = inject(ActivatedRoute);
  private postRepo = inject(POST_REPOSITORY);
  public authService = inject(AuthService);
  private cdr = inject(ChangeDetectorRef);

  ngOnInit(): void {
    const slug = this.route.snapshot.paramMap.get('slug') || '';
    this.postRepo.getPostBySlug(slug).subscribe({
      next: (post) => {
        this.post = post || null;
        this.isLoading = false;
        if (!post) {
          this.error = 'Post no encontrado.';
        }
        this.cdr.markForCheck();
      },
      error: () => {
        this.error = 'Error al cargar el post.';
        this.isLoading = false;
        this.cdr.markForCheck();
      }
    });
  }

  async submitComment() {
    if (!this.newCommentText.trim() || !this.post) return;
    
    this.isSubmittingComment = true;
    this.commentError = '';
    
    const userId = await this.authService.getCurrentUserId();
    if (!userId) {
      this.commentError = 'Debes iniciar sesión para comentar.';
      this.isSubmittingComment = false;
      this.cdr.markForCheck();
      return;
    }

    const commentPayload = {
      id: '',
      author: userId,
      text: this.newCommentText,
      date: new Date()
    };

    this.postRepo.addComment(this.post.id, commentPayload).subscribe({
      next: (updatedPost) => {
        this.post = updatedPost;
        this.newCommentText = '';
        this.isSubmittingComment = false;
        this.cdr.markForCheck();
      },
      error: (err) => {
        this.commentError = 'Error al enviar el comentario.';
        this.isSubmittingComment = false;
        this.cdr.markForCheck();
      }
    });
  }
}
