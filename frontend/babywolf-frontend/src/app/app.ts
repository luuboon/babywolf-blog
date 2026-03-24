import { Component } from '@angular/core';
import { RouterModule } from '@angular/router';

import { Header } from './core/layout/header/header';
import { ExternalLinks } from './shared/ui/external-links/external-links';

@Component({
  selector: 'app-root',
  imports: [RouterModule, Header, ExternalLinks],
  templateUrl: './app.html',
  styleUrls: ['./app.scss'],
})
export class App { }
