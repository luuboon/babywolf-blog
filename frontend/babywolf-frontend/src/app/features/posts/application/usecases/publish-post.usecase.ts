import { Observable } from 'rxjs';
import { UseCase } from '../../../../core/application/base/use-case';
import { Post } from '../../domain/models/post.model';
import { PostRepository } from '../../domain/repositories/post.repository';

// Assuming publishing means updating a draft to published, or just creating it.
// Given strict clean architecture, this might encapsulate specific business rules for publishing.
export class PublishPost implements UseCase<string, Post> {
    constructor(private postRepository: PostRepository) { }

    execute(postId: string): Observable<Post> {
        // In a real app, this would set status = 'PUBLISHED'
        return this.postRepository.updatePost(postId, {});
    }
}
