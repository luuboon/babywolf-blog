import { Observable } from 'rxjs';
import { Post, Comment } from '../models/post.model';

export interface PostRepository {
    getPosts(): Observable<Post[]>;
    getPostById(id: string): Observable<Post | undefined>;
    getPostBySlug(slug: string): Observable<Post | undefined>;
    createPost(post: Omit<Post, 'id' | 'date'>): Observable<Post>;
    updatePost(id: string, post: Partial<Post>): Observable<Post>;
    deletePost(id: string): Observable<void>;
    addComment(postId: string, comment: Comment): Observable<Post>;
}
