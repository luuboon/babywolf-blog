import { Component, Input } from '@angular/core';
import { DatePipe } from '@angular/common';
import { RouterModule } from '@angular/router';
import { Post } from '../../../domain/models/post.model';

@Component({
    selector: 'app-post-card',
    standalone: true,
    imports: [DatePipe, RouterModule],
    templateUrl: './post-card.html',
    styleUrls: ['./post-card.scss']
})
export class PostCard {
    @Input({ required: true }) post!: Post;

    // Practica 9-10: "destacar contenido mediante el puntero" con eventos DOM explícitos.
    highlighted = false;

    onPointerEnter(): void {
        this.highlighted = true;
    }

    onPointerLeave(): void {
        this.highlighted = false;
    }
}
