import { Injectable, inject } from '@angular/core';
import { AuthService } from './auth.service';

const IDLE_LIMIT_MS = 20 * 60 * 1000; // 20 min de inactividad
const ACTIVITY_EVENTS = ['click', 'keydown', 'mousemove', 'scroll', 'touchstart'] as const;

/**
 * Practica 11-12, control de seguridad "expiración de sesión": cierra la
 * sesión automáticamente si el usuario no interactúa con la página.
 * Se suma a la expiración propia del JWT de Supabase (exp del token).
 */
@Injectable({ providedIn: 'root' })
export class SessionTimeoutService {
  private auth = inject(AuthService);
  private timer?: ReturnType<typeof setTimeout>;
  private started = false;

  start(): void {
    if (this.started || typeof window === 'undefined') return;
    this.started = true;

    ACTIVITY_EVENTS.forEach(evt =>
      window.addEventListener(evt, () => this.resetTimer(), { passive: true })
    );
    this.resetTimer();
  }

  private resetTimer(): void {
    if (this.timer) clearTimeout(this.timer);
    this.timer = setTimeout(() => this.auth.signOut(), IDLE_LIMIT_MS);
  }
}
