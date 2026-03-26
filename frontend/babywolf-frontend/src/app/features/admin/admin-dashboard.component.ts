import { Component, OnInit, inject, ChangeDetectorRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { SupabaseService } from '../../core/services/supabase.service';

@Component({
  selector: 'app-admin-dashboard',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './admin-dashboard.component.html',
  styleUrls: ['./admin-dashboard.component.scss']
})
export class AdminDashboardComponent implements OnInit {
  
  stats = [
    { title: 'Usuarios', value: 0, icon: '👥', gradient: 'grad-blue' },
    { title: 'Posts', value: 0, icon: '📝', gradient: 'grad-red' },
    { title: 'Comentarios', value: 0, icon: '💬', gradient: 'grad-purple' },
    { title: 'Publicados', value: 0, icon: '✅', gradient: 'grad-green' }
  ];

  quickActions = [
    { title: 'Gestión de Usuarios', desc: 'Ver, promover o suspender usuarios registrados.', icon: '👥', link: '/admin/users', btnText: 'Ir a Usuarios' },
    { title: 'Gestión de Posts', desc: 'Crear, editar o eliminar artículos del blog.', icon: '📝', link: '/admin/posts', btnText: 'Ir a Posts' },
    { title: 'Moderación', desc: 'Revisar y moderar comentarios de los lectores.', icon: '💬', link: '/admin/comments', btnText: 'Ir a Comentarios' }
  ];

  loading = true;

  private sb = inject(SupabaseService);
  private cdr = inject(ChangeDetectorRef);

  ngOnInit() {
    this.loadStats();
  }

  async loadStats() {
    this.loading = true;

    const [usersRes, postsRes, commentsRes, publishedRes] = await Promise.all([
      this.sb.client.from('users').select('id', { count: 'exact', head: true }),
      this.sb.client.from('posts').select('id', { count: 'exact', head: true }),
      this.sb.client.from('comments').select('id', { count: 'exact', head: true }),
      this.sb.client.from('posts').select('id', { count: 'exact', head: true }).eq('published', true)
    ]);

    this.stats[0].value = usersRes.count ?? 0;
    this.stats[1].value = postsRes.count ?? 0;
    this.stats[2].value = commentsRes.count ?? 0;
    this.stats[3].value = publishedRes.count ?? 0;

    this.loading = false;
    this.cdr.markForCheck();
  }
}
