import { Injectable, inject } from '@angular/core';
import { Observable, from, of, throwError } from 'rxjs';
import { map, switchMap } from 'rxjs/operators';
import { Post, Comment } from '../../domain/models/post.model';
import { PostRepository } from '../../domain/repositories/post.repository';
import { SupabaseService } from '../../../../core/services/supabase.service';

/**
 * DTO that mirrors the Supabase `posts` table + joined `users` data.
 * snake_case matches the DB column names.
 */
interface PostRow {
    id: string;
    author_id: string;
    title: string;
    slug: string;
    content: string;
    cover_image_url: string | null;
    published: boolean;
    category: string;
    created_at: string;
    updated_at: string;
    users?: { email: string; username: string } | null;
}

interface CommentRow {
    id: string;
    post_id: string;
    author_id: string;
    content: string;
    created_at: string;
    users?: { username: string } | null;
}

@Injectable({
    providedIn: 'root',
})
export class SupabasePostRepository implements PostRepository {
    private sb = inject(SupabaseService);

    /**
     * Map a DB row (snake_case) → Domain Model (camelCase).
     */
    private toDomain(row: PostRow, comments: Comment[] = []): Post {
        return {
            id: row.id,
            title: row.title,
            slug: row.slug,
            excerpt: row.content ? row.content.substring(0, 160) + '...' : '',
            content: row.content,
            author: row.users?.username || row.users?.email || 'Unknown',
            authorId: row.author_id,
            date: new Date(row.created_at),
            category: row.category || 'general',
            coverImageUrl: row.cover_image_url || '',
            published: row.published,
            comments,
            likes: 0 // No likes column in DB yet
        };
    }

    private commentToDomain(row: CommentRow): Comment {
        return {
            id: row.id,
            author: row.users?.username || 'Anonymous',
            text: row.content,
            date: new Date(row.created_at)
        };
    }

    getPosts(): Observable<Post[]> {
        return from(
            this.sb.client
                .from('posts')
                .select('*, users(email, username)')
                .eq('published', true)
                .order('created_at', { ascending: false })
        ).pipe(
            map(({ data, error }) => {
                if (error) throw error;
                return (data as PostRow[]).map(row => this.toDomain(row));
            })
        );
    }

    getPostById(id: string): Observable<Post | undefined> {
        return from(
            this.sb.client
                .from('posts')
                .select('*, users(email, username)')
                .eq('id', id)
                .single()
        ).pipe(
            switchMap(({ data, error }) => {
                if (error) return of(undefined);
                const post = data as PostRow;
                return this.getCommentsForPost(post.id).pipe(
                    map(comments => this.toDomain(post, comments))
                );
            })
        );
    }

    getPostBySlug(slug: string): Observable<Post | undefined> {
        return from(
            this.sb.client
                .from('posts')
                .select('*, users(email, username)')
                .eq('slug', slug)
                .single()
        ).pipe(
            switchMap(({ data, error }) => {
                if (error) return of(undefined);
                const post = data as PostRow;
                return this.getCommentsForPost(post.id).pipe(
                    map(comments => this.toDomain(post, comments))
                );
            })
        );
    }

    createPost(post: Omit<Post, 'id' | 'date'>): Observable<Post> {
        const row = {
            title: post.title,
            slug: post.slug,
            content: post.content,
            cover_image_url: post.coverImageUrl || null,
            published: post.published,
            category: post.category || 'general',
            author_id: post.authorId
        };

        return from(
            this.sb.client
                .from('posts')
                .insert(row)
                .select('*, users(email, username)')
                .single()
        ).pipe(
            map(({ data, error }) => {
                if (error) throw error;
                return this.toDomain(data as PostRow);
            })
        );
    }

    updatePost(id: string, post: Partial<Post>): Observable<Post> {
        const row: Record<string, unknown> = {};
        if (post.title !== undefined) row['title'] = post.title;
        if (post.slug !== undefined) row['slug'] = post.slug;
        if (post.content !== undefined) row['content'] = post.content;
        if (post.coverImageUrl !== undefined) row['cover_image_url'] = post.coverImageUrl;
        if (post.published !== undefined) row['published'] = post.published;
        if (post.category !== undefined) row['category'] = post.category;

        return from(
            this.sb.client
                .from('posts')
                .update(row)
                .eq('id', id)
                .select('*, users(email, username)')
                .single()
        ).pipe(
            map(({ data, error }) => {
                if (error) throw error;
                return this.toDomain(data as PostRow);
            })
        );
    }

    deletePost(id: string): Observable<void> {
        return from(
            this.sb.client
                .from('posts')
                .delete()
                .eq('id', id)
        ).pipe(
            map(({ error }) => {
                if (error) throw error;
            })
        );
    }

    addComment(postId: string, comment: Comment): Observable<Post> {
        return from(
            this.sb.client
                .from('comments')
                .insert({
                    post_id: postId,
                    author_id: comment.author, // Will need to be user ID
                    content: comment.text
                })
        ).pipe(
            switchMap(({ error }) => {
                if (error) return throwError(() => error);
                return this.getPostById(postId).pipe(
                    map(post => {
                        if (!post) throw new Error('Post not found after comment insert');
                        return post;
                    })
                );
            })
        );
    }

    /**
     * Fetch comments for a specific post.
     */
    private getCommentsForPost(postId: string): Observable<Comment[]> {
        return from(
            this.sb.client
                .from('comments')
                .select('*, users(username)')
                .eq('post_id', postId)
                .order('created_at', { ascending: true })
        ).pipe(
            map(({ data, error }) => {
                if (error || !data) return [];
                return (data as CommentRow[]).map(row => this.commentToDomain(row));
            })
        );
    }
}
