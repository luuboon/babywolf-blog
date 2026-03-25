import { Component, inject } from '@angular/core';
import { CommonModule } from '@angular/common';
import { Router, RouterModule, RouterOutlet } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { AuthService } from '../../../core/services/auth.service';

interface MenuItem {
  label: string;
  icon: string;
  routerLink?: string;
}

@Component({
  selector: 'app-user-layout',
  standalone: true,
  imports: [CommonModule, RouterModule, RouterOutlet, FormsModule],
  templateUrl: './user-layout.html',
  styleUrls: ['./user-layout.scss']
})
export class UserLayoutComponent {
  isSidebarCollapsed = false;
  searchQuery = '';
  
  authService = inject(AuthService);
  currentUser$ = this.authService.currentUser$;
  isAdmin$ = this.authService.isAdmin$;

  mainMenu: MenuItem[] = [
    { label: 'Inicio', icon: '🏠', routerLink: '/' },
    { label: 'Posts', icon: '📝', routerLink: '/posts' },
    { label: 'Retro', icon: '🕹️', routerLink: '/category/retro' },
    { label: 'Gaming', icon: '🎮', routerLink: '/category/gaming' },
    { label: 'Opinión', icon: '💭', routerLink: '/category/opinion' },
    { label: 'Tech', icon: '💻', routerLink: '/category/tech' },
  ];

  constructor(private readonly router: Router) {}

  toggleSidebar(): void {
    this.isSidebarCollapsed = !this.isSidebarCollapsed;
  }

  onSearch(): void {
    if (this.searchQuery.trim()) {
      this.router.navigate(['/search'], { queryParams: { q: this.searchQuery } });
      this.searchQuery = '';
    }
  }

  onSearchKeydown(event: KeyboardEvent): void {
    if (event.key === 'Enter') {
      this.onSearch();
    }
  }

  async logout() {
    await this.authService.signOut();
    this.router.navigate(['/']);
  }
}
