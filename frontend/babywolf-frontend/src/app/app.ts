import { Component, inject } from '@angular/core';
import { RouterOutlet } from '@angular/router';
import { SessionTimeoutService } from './core/services/session-timeout.service';
import { SeasonalThemeService } from './core/services/seasonal-theme.service';

@Component({
  selector: 'app-root',
  imports: [RouterOutlet],
  template: '<router-outlet></router-outlet>',
  styles: [':host { display: block; min-height: 100vh; }'],
})
export class App {
  private sessionTimeout = inject(SessionTimeoutService);
  private seasonalTheme = inject(SeasonalThemeService);

  constructor() {
    this.sessionTimeout.start();
    this.seasonalTheme.apply();
  }
}
