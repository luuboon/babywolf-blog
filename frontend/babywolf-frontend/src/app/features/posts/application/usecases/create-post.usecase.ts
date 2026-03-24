import { Observable } from 'rxjs';
import { UseCase } from '../../../../core/application/base/use-case';
import { Post } from '../../domain/models/post.model';
import { PostRepository } from '../../domain/repositories/post.repository';

export class CreatePost implements UseCase<Omit<Post, 'id' | 'date'>, Post> {
    constructor(private postRepository: PostRepository) { }

    execute(params: Omit<Post, 'id' | 'date'>): Observable<Post> {
        return this.postRepository.createPost(params);
    }
}
