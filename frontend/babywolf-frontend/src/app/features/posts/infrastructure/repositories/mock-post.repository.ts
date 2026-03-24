import { Injectable } from '@angular/core';
import { Observable, of, throwError } from 'rxjs';
import { Post, Comment } from '../../domain/models/post.model';
import { PostRepository } from '../../domain/repositories/post.repository';

@Injectable({
    providedIn: 'root',
})
export class MockPostRepository implements PostRepository {
    private posts: Post[] = [
        {
            id: '1',
            title: 'Why the Dreamcast Deserved Better',
            slug: 'why-the-dreamcast-deserved-better',
            excerpt: 'Evaluating the innovative features of Sega\'s final console.',
            content: 'The Dreamcast was ahead of its time. With built-in internet capabilities, the VMU, and an incredible library of arcade ports, it set the stage for modern online gaming. However, poor timing and the looming PS2 launch cut its life short.',
            author: 'RetroWolf',
            authorId: 'u1',
            date: new Date('2023-10-27'),
            category: 'Retro',
            coverImageUrl: '',
            published: true,
            comments: [
                { id: 'c1', author: 'SonicFan', text: 'It still looks good on VGA!', date: new Date('2023-10-28') }
            ],
            likes: 128
        },
        {
            id: '2',
            title: 'SNES vs Genesis: The Sound Chip War',
            slug: 'snes-vs-genesis-the-sound-chip-war',
            excerpt: 'Orchestral samples vs FM Synthesis. Which one won?',
            content: 'The SNES used the SPC700 for sample-based audio, allowing for realistic instrument sounds. The Genesis used the YM2612 FM chip, producing gritty, synthesized basslines. Each had its strengths, defining the "sound" of the 16-bit era.',
            author: 'AudioBit',
            authorId: 'u2',
            date: new Date('2023-11-05'),
            category: 'Opinión',
            coverImageUrl: '',
            published: true,
            comments: [
                { id: 'c2', author: 'BlastProcessing', text: 'Genesis does what Nintendon\'t.', date: new Date('2023-11-06') }
            ],
            likes: 84
        },
        {
            id: '3',
            title: 'The Art of PS1 Polygons',
            slug: 'the-art-of-ps1-polygons',
            excerpt: 'Why low-poly aesthetics are making a comeback.',
            content: 'Wobbly textures and low polygon counts were once seen as technical limitations. Today, they are an aesthetic choice. The PS1 era has a unique charm that modern photorealism simply cannot replicate.',
            author: 'PolyGone',
            authorId: 'u3',
            date: new Date('2023-11-15'),
            category: 'Gaming',
            coverImageUrl: '',
            published: true,
            comments: [],
            likes: 256
        }
    ];

    getPosts(): Observable<Post[]> {
        return of(this.posts);
    }

    getPostById(id: string): Observable<Post | undefined> {
        const post = this.posts.find(p => p.id === id);
        return of(post);
    }

    getPostBySlug(slug: string): Observable<Post | undefined> {
        const post = this.posts.find(p => p.slug === slug);
        return of(post);
    }

    createPost(post: Omit<Post, 'id' | 'date'>): Observable<Post> {
        const newPost: Post = {
            ...post,
            id: (this.posts.length + 1).toString(),
            date: new Date()
        };
        this.posts.push(newPost);
        return of(newPost);
    }

    updatePost(id: string, post: Partial<Post>): Observable<Post> {
        const index = this.posts.findIndex(p => p.id === id);
        if (index !== -1) {
            this.posts[index] = { ...this.posts[index], ...post };
            return of(this.posts[index]);
        }
        return throwError(() => new Error('Post not found'));
    }

    deletePost(id: string): Observable<void> {
        this.posts = this.posts.filter(p => p.id !== id);
        return of(undefined);
    }

    addComment(postId: string, comment: Comment): Observable<Post> {
        const index = this.posts.findIndex(p => p.id === postId);
        if (index !== -1) {
            this.posts[index].comments.push(comment);
            return of(this.posts[index]);
        }
        return throwError(() => new Error('Post not found'));
    }
}
