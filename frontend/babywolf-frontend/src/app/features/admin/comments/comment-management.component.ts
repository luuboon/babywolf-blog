import { Component, OnInit, inject, ChangeDetectorRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { SupabaseService } from '../../../core/services/supabase.service';

interface Comment {
  id: string;
  content: string;
  created_at: string;
  users?: { email: string; username: string } | null;
  posts?: { title: string } | null;
}

@Component({
  selector: 'app-comment-management',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './comment-management.component.html',
  styleUrls: ['./comment-management.component.scss']
})
export class CommentManagementComponent implements OnInit {
  comments: Comment[] = [];
  loading = true;
  errorMsg = '';

  private sb = inject(SupabaseService);
  private cdr = inject(ChangeDetectorRef);

  ngOnInit() {
    this.loadComments();
  }

  async loadComments() {
    this.loading = true;
    const { data, error } = await this.sb.client
      .from('comments')
      .select('*, users(email, username), posts(title)')
      .order('created_at', { ascending: false });

    if (error) {
      this.errorMsg = 'Error al cargar los comentarios: ' + error.message;
      console.error(error);
    } else {
      this.comments = data as Comment[];
    }
    this.loading = false;
    this.cdr.markForCheck();
  }

  async deleteComment(id: string) {
    if (!confirm('¿Seguro de que deseas eliminar este comentario? Acción irreversible.')) return;
    
    const { error } = await this.sb.client
      .from('comments')
      .delete()
      .eq('id', id);

    if (error) {
      alert('Error al eliminar comentario: ' + error.message);
      console.error(error);
    } else {
      this.comments = this.comments.filter(c => c.id !== id);
    }
    this.cdr.markForCheck();
  }
}
