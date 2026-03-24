import { Observable } from 'rxjs';
import { UseCase } from '../../../../core/application/base/use-case';
import { Post } from '../../domain/models/post.model';
import { PostRepository } from '../../domain/repositories/post.repository';

export class GetPostBySlug implements UseCase<string, Post | undefined> {
    constructor(private postRepository: PostRepository) { }

    execute(slug: string): Observable<Post | undefined> {
        return this.postRepository.getPostBySlug(slug);
    }
}
