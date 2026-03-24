export interface Comment {
    id: string;
    author: string;
    text: string;
    date: Date;
}

export interface Post {
    id: string;
    title: string;
    slug: string;
    excerpt: string;
    content: string;
    author: string;
    authorId: string;
    date: Date;
    category: string;
    coverImageUrl: string;
    published: boolean;
    comments: Comment[];
    likes: number;
}
