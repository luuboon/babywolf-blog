import { Component, inject } from '@angular/core';
import { RouterLink } from '@angular/router';

import { MatToolbarModule } from '@angular/material/toolbar';
import { MatButtonModule } from '@angular/material/button';
import { SearchStateService } from '../../services/search-state.service';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [RouterLink, MatToolbarModule, MatButtonModule],
  templateUrl: './header.html',
  styleUrls: ['./header.scss'],
})
export class Header {
  private searchService = inject(SearchStateService);

  onSearch(event: Event): void {
    const input = event.target as HTMLInputElement;
    this.searchService.setSearchQuery(input.value);
  }
}

