import { Observable } from 'rxjs';
import { UseCase } from '../../../../core/application/base/use-case';
import { Post, Comment } from '../../domain/models/post.model';
import { PostRepository } from '../../domain/repositories/post.repository';

export class AddComment implements UseCase<{ postId: string, comment: Comment }, Post> {
    constructor(private postRepository: PostRepository) { }

    execute(params: { postId: string, comment: Comment }): Observable<Post> {
        return this.postRepository.addComment(params.postId, params.comment);
    }
}
