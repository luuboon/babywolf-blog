import { Component, OnInit, inject, ChangeDetectorRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, RouterModule } from '@angular/router';
import { Post } from '../posts/domain/models/post.model';
import { POST_REPOSITORY } from '../posts/domain/repositories/post.repository.token';

@Component({
  selector: 'app-category-posts',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './category-posts.html',
  styleUrls: ['./category-posts.scss']
})
export class CategoryPostsPage implements OnInit {
  categoryName = '';
  posts: Post[] = [];
  isLoading = true;

  private route = inject(ActivatedRoute);
  private postRepo = inject(POST_REPOSITORY);
  private cdr = inject(ChangeDetectorRef);

  ngOnInit(): void {
    this.route.paramMap.subscribe(params => {
      this.categoryName = params.get('name') || '';
      this.loadPosts();
    });
  }

  private loadPosts(): void {
    this.isLoading = true;
    this.postRepo.getPosts().subscribe({
      next: (posts) => {
        this.posts = posts.filter(p =>
          p.category.toLowerCase() === this.categoryName.toLowerCase()
        );
        this.isLoading = false;
        this.cdr.markForCheck();
      },
      error: () => {
        this.isLoading = false;
        this.cdr.markForCheck();
      }
    });
  }
}
