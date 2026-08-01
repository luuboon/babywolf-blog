import { Component, ElementRef, HostListener, inject } from '@angular/core';
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

  // Practica 7: menú desplegable disparado por un evento de puntero (click).
  categoriesMenuOpen = false;
  readonly categoriesMenu = [
    { label: '🎮 Gaming', routerLink: '/category/gaming' },
    { label: '💻 Tech', routerLink: '/category/tech' },
    { label: '💭 Opinión', routerLink: '/category/opinion' },
    { label: '🕹️ Retro', routerLink: '/category/retro' },
  ];

  constructor(private readonly router: Router, private readonly elementRef: ElementRef<HTMLElement>) {}

  toggleSidebar(): void {
    this.isSidebarCollapsed = !this.isSidebarCollapsed;
  }

  toggleCategoriesMenu(event: MouseEvent): void {
    event.stopPropagation();
    this.categoriesMenuOpen = !this.categoriesMenuOpen;
  }

  closeCategoriesMenu(): void {
    this.categoriesMenuOpen = false;
  }

  @HostListener('document:click', ['$event'])
  onDocumentClick(event: MouseEvent): void {
    if (this.categoriesMenuOpen && !this.elementRef.nativeElement.contains(event.target as Node)) {
      this.categoriesMenuOpen = false;
    }
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
