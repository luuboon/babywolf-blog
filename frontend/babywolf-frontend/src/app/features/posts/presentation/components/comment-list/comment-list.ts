import { Component, Input } from '@angular/core';
import { DatePipe } from '@angular/common';
import { Comment } from '../../../domain/models/post.model';

@Component({
    selector: 'app-comment-list',
    standalone: true,
    imports: [DatePipe],
    templateUrl: './comment-list.html',
    styleUrls: ['./comment-list.scss']
})
export class CommentList {
    @Input() comments: Comment[] = [];
}
