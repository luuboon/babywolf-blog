import { Component } from '@angular/core';
import { RouterOutlet } from '@angular/router';

import { Header } from './shared/layout/header/header';
import { ExternalLinks } from './shared/external-links/external-links';

@Component({
  selector: 'app-root',
  imports: [RouterOutlet, Header, ExternalLinks],
  templateUrl: './app.html',
  styleUrls: ['./app.scss'],
})
export class App {}
