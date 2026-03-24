import { Component, OnInit, inject } from '@angular/core';
import { CommonModule } from '@angular/common';

import { PostCard } from '../../components/post-card/post-card';
import { PostStateService } from '../../state/post-state.service';
import { POST_REPOSITORY } from '../../../domain/repositories/post.repository.token';

@Component({
    selector: 'app-posts-page',
    standalone: true,
    imports: [CommonModule, PostCard],
    templateUrl: './posts-page.html',
    styleUrls: ['./posts-page.scss'],
    providers: [
        PostStateService
    ]
})
export class PostsPage implements OnInit {
    protected state = inject(PostStateService);

    ngOnInit(): void {
        this.state.loadPosts();
    }
}

