import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule, RouterOutlet } from '@angular/router';

interface AdminMenuItem {
  label: string;
  icon: string;
  routerLink: string;
}

@Component({
  selector: 'app-admin-layout',
  standalone: true,
  imports: [CommonModule, RouterModule, RouterOutlet],
  templateUrl: './admin-layout.html',
  styleUrls: ['./admin-layout.scss']
})
export class AdminLayoutComponent {
  isSidebarCollapsed = false;

  adminMenu: AdminMenuItem[] = [
    { label: 'Dashboard', icon: '📊', routerLink: '/admin' },
    { label: 'Usuarios', icon: '👥', routerLink: '/admin/users' },
    { label: 'Posts', icon: '📝', routerLink: '/admin/posts' },
    { label: 'Comentarios', icon: '💬', routerLink: '/admin/comments' },
  ];

  toggleSidebar(): void {
    this.isSidebarCollapsed = !this.isSidebarCollapsed;
  }
}
