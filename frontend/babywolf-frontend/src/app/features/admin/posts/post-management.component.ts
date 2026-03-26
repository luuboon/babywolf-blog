import { Component, OnInit, inject, ChangeDetectorRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { SupabaseService } from '../../../core/services/supabase.service';

interface Post {
  id: string;
  title: string;
  slug: string;
  category: string;
  created_at: string;
  published: boolean;
  users?: { email: string; username: string } | null;
}

@Component({
  selector: 'app-post-management',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './post-management.component.html',
  styleUrls: ['../admin-dashboard.component.scss']
})
export class PostManagementComponent implements OnInit {
  posts: Post[] = [];
  loading = true;
  errorMsg = '';

  private sb = inject(SupabaseService);
  private cdr = inject(ChangeDetectorRef);

  ngOnInit() {
    this.loadPosts();
  }

  async loadPosts() {
    this.loading = true;
    const { data, error } = await this.sb.client
      .from('posts')
      .select('*, users(email, username)')
      .order('created_at', { ascending: false });

    if (error) {
      this.errorMsg = 'Error al cargar posts: ' + error.message;
      console.error(error);
    } else {
      this.posts = data as Post[];
    }
    this.loading = false;
    this.cdr.markForCheck();
  }

  async deletePost(id: string) {
    if (!confirm('¿Seguro de que deseas eliminar este post permanentemente?')) return;
    
    const { error } = await this.sb.client
      .from('posts')
      .delete()
      .eq('id', id);

    if (error) {
      alert('Error al eliminar post: ' + error.message);
      console.error(error);
    } else {
      this.posts = this.posts.filter(p => p.id !== id);
    }
    this.cdr.markForCheck();
  }

  async togglePublish(post: Post) {
    const { error } = await this.sb.client
      .from('posts')
      .update({ published: !post.published })
      .eq('id', post.id);

    if (error) {
      alert('Error al actualizar estado de publicación: ' + error.message);
      console.error(error);
    } else {
      post.published = !post.published;
    }
    this.cdr.markForCheck();
  }
}
