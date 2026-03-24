import { Component, inject, OnInit } from '@angular/core';
import { CommonModule, DatePipe } from '@angular/common';
import { RouterModule } from '@angular/router';
import { Post } from '../../../posts/domain/models/post.model';
import { POST_REPOSITORY } from '../../../posts/domain/repositories/post.repository.token';

@Component({
  selector: 'app-home',
  standalone: true,
  imports: [CommonModule, RouterModule, DatePipe],
  templateUrl: './home.component.html',
  styleUrls: ['./home.component.scss']
})
export class HomeComponent implements OnInit {
  recentPosts: Post[] = [];
  isLoading = true;

  private postRepo = inject(POST_REPOSITORY);

  ngOnInit(): void {
    this.postRepo.getPosts().subscribe({
      next: (posts) => {
        this.recentPosts = posts.slice(0, 6);
        this.isLoading = false;
      },
      error: () => {
        this.isLoading = false;
      }
    });
  }
}
