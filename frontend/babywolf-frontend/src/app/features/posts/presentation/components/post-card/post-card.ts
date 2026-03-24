import { Component, Input } from '@angular/core';
import { DatePipe } from '@angular/common';
import { Post } from '../../../domain/models/post.model';

@Component({
    selector: 'app-post-card',
    standalone: true,
    imports: [DatePipe],
    templateUrl: './post-card.html',
    styleUrls: ['./post-card.scss']
})
export class PostCard {
    @Input({ required: true }) post!: Post;
}
