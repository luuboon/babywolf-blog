import { Component, OnInit, inject, ChangeDetectorRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ActivatedRoute, RouterModule } from '@angular/router';
import { Post } from '../posts/domain/models/post.model';
import { POST_REPOSITORY } from '../posts/domain/repositories/post.repository.token';

@Component({
  selector: 'app-search-results',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './search-results.html',
  styleUrls: ['./search-results.scss']
})
export class SearchResultsPage implements OnInit {
  query = '';
  results: Post[] = [];
  isLoading = true;

  private route = inject(ActivatedRoute);
  private postRepo = inject(POST_REPOSITORY);
  private cdr = inject(ChangeDetectorRef);

  ngOnInit(): void {
    this.route.queryParams.subscribe(params => {
      this.query = params['q'] || '';
      this.search();
    });
  }

  private search(): void {
    if (!this.query.trim()) {
      this.results = [];
      this.isLoading = false;
      this.cdr.markForCheck();
      return;
    }

    this.isLoading = true;
    this.postRepo.getPosts().subscribe({
      next: (posts) => {
        const q = this.query.toLowerCase();
        this.results = posts.filter(p =>
          p.title.toLowerCase().includes(q) ||
          p.excerpt.toLowerCase().includes(q) ||
          p.content.toLowerCase().includes(q) ||
          p.category.toLowerCase().includes(q)
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
