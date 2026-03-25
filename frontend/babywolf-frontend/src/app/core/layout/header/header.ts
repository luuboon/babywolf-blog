import { Component, inject } from '@angular/core';
import { Router, RouterLink } from '@angular/router';
import { toSignal } from '@angular/core/rxjs-interop';

import { MatToolbarModule } from '@angular/material/toolbar';
import { MatButtonModule } from '@angular/material/button';
import { SearchStateService } from '../../services/search-state.service';
import { AuthService } from '../../services/auth.service';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [RouterLink, MatToolbarModule, MatButtonModule],
  templateUrl: './header.html',
  styleUrls: ['./header.scss'],
})
export class Header {
  private searchService = inject(SearchStateService);
  private authService = inject(AuthService);
  private router = inject(Router);

  currentUser = toSignal(this.authService.currentUser$);
  isAdmin = toSignal(this.authService.isAdmin$);

  onSearch(event: Event): void {
    const input = event.target as HTMLInputElement;
    this.searchService.setSearchQuery(input.value);
  }

  async logout() {
    await this.authService.signOut();
    this.router.navigate(['/']);
  }
}


