import { Component, OnInit, inject } from '@angular/core';
import { CommonModule } from '@angular/common'; // For AsyncPipe if needed, but we use signals

import { PostCard } from '../../components/post-card/post-card';
import { PostStateService } from '../../state/post-state.service';
import { MockPostRepository } from '../../../infrastructure/repositories/mock-post.repository';
import { POST_REPOSITORY } from '../../../domain/repositories/post.repository.token';

@Component({
    selector: 'app-posts-page',
    standalone: true,
    imports: [CommonModule, PostCard],
    templateUrl: './posts-page.html',
    styleUrls: ['./posts-page.scss'],
    providers: [
        PostStateService,
        { provide: POST_REPOSITORY, useClass: MockPostRepository }
    ]
})
export class PostsPage implements OnInit {
    // Inject the state service (it's provided locally, so it's a fresh instance for this page tree)
    // In a real app, maybe we want it global, but per-feature is good too.
    protected state = inject(PostStateService);

    ngOnInit(): void {
        this.state.loadPosts();
    }
}
