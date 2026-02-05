import { Component } from '@angular/core';
import { RouterLink } from '@angular/router';

import { MatToolbarModule } from '@angular/material/toolbar';
import { MatButtonModule } from '@angular/material/button';

import { ExternalLinks } from '../../external-links/external-links';

@Component({
  selector: 'app-header',
  standalone: true,
  imports: [RouterLink, MatToolbarModule, MatButtonModule, ExternalLinks],
  templateUrl: './header.html',
  styleUrls: ['./header.scss'],
})
export class Header { }

