import { Injectable, inject, signal, computed } from '@angular/core';
import { Post } from '../../domain/models/post.model';
import { POST_REPOSITORY } from '../../domain/repositories/post.repository.token';
import { finalize } from 'rxjs/operators';

@Injectable()
export class PostStateService {
    private repository = inject(POST_REPOSITORY);

    // State Signals
    private postsSignal = signal<Post[]>([]);
    private loadingSignal = signal<boolean>(false);
    private errorSignal = signal<string | null>(null);

    // Read-only Signals
    readonly posts = computed(() => this.postsSignal());
    readonly isLoading = computed(() => this.loadingSignal());
    readonly error = computed(() => this.errorSignal());

    loadPosts(): void {
        if (this.postsSignal().length > 0) return; // Prevent unnecessary reloads

        this.loadingSignal.set(true);
        this.errorSignal.set(null);

        this.repository.getPosts()
            .pipe(finalize(() => this.loadingSignal.set(false)))
            .subscribe({
                next: (posts) => this.postsSignal.set(posts),
                error: (err) => this.errorSignal.set('Failed to load posts. Please try again later.')
            });
    }
}
