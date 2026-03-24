export interface Comment {
    id: string;
    author: string;
    text: string;
    date: Date;
}

export interface Post {
    id: string;
    title: string;
    excerpt: string;
    content: string;
    author: string;
    date: Date;
    consoleCategory: 'SNES' | 'SEGA' | 'PS1' | 'N64' | 'DREAMCAST' | 'OTHER';
    comments: Comment[];
    likes: number;
}
