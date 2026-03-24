import { InjectionToken } from '@angular/core';
import { PostRepository } from './post.repository';

export const POST_REPOSITORY = new InjectionToken<PostRepository>('PostRepository');
