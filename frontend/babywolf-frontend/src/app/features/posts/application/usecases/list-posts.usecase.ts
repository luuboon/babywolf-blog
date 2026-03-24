import { Observable } from 'rxjs';
import { UseCase } from '../../../../core/application/base/use-case';
import { Post } from '../../domain/models/post.model';
import { PostRepository } from '../../domain/repositories/post.repository';

export class ListPosts implements UseCase<void, Post[]> {
    constructor(private postRepository: PostRepository) { }

    execute(): Observable<Post[]> {
        return this.postRepository.getPosts();
    }
}
