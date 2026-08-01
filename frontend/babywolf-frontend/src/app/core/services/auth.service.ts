import { Injectable, inject } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { User } from '@supabase/supabase-js';
import { SupabaseService } from './supabase.service';

/**
 * Política de contraseña segura (Practica 11-12, parte 4):
 * min. 8 caracteres, mayúscula, minúscula, número y símbolo.
 */
export const STRONG_PASSWORD_PATTERN =
  /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[^A-Za-z0-9]).{8,}$/;
export const STRONG_PASSWORD_HINT =
  'Mínimo 8 caracteres, con mayúscula, minúscula, número y símbolo.';

@Injectable({
  providedIn: 'root'
})
export class AuthService {
  private sb = inject(SupabaseService);
  private currentUserSubject = new BehaviorSubject<User | null>(null);
  public currentUser$: Observable<User | null> = this.currentUserSubject.asObservable();
  
  private isAdminSubject = new BehaviorSubject<boolean>(false);
  public isAdmin$: Observable<boolean> = this.isAdminSubject.asObservable();

  constructor() {
    // Check active session on init
    this.sb.client.auth.getSession().then(({ data }) => {
       if (data.session) {
         this.currentUserSubject.next(data.session.user);
         this.loadUserRole(data.session.user.id);
       }
    });

    this.sb.client.auth.onAuthStateChange((event, session) => {
      if (session) {
        this.currentUserSubject.next(session.user);
        this.loadUserRole(session.user.id);

        if (event === 'SIGNED_IN') {
          // Enviar email de notificación de login via Resend
          fetch('/api/notify-login', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({
              email: session.user.email,
              timestamp: new Date().toISOString()
            })
          })
          .then(r => r.json())
          .then(d => console.log('[Seguridad] Email de login enviado:', d))
          .catch(e => console.warn('[Seguridad] No se pudo enviar email:', e));
        }
      } else {
        this.currentUserSubject.next(null);
        this.isAdminSubject.next(false);
      }
    });
  }

  private async loadUserRole(userId: string) {
    const { data } = await this.sb.client.from('users').select('role').eq('id', userId).single();
    this.isAdminSubject.next(data?.role === 'admin');
  }

  async checkIsAdmin(): Promise<boolean> {
    const { data: { session } } = await this.sb.client.auth.getSession();
    if (!session) return false;
    const { data } = await this.sb.client.from('users').select('role').eq('id', session.user.id).single();
    return data?.role === 'admin';
  }

  async signUp(email: string, password: string) {
    return this.sb.client.auth.signUp({ email, password });
  }

  /**
   * Login con verificación de bloqueo temporal (Practica 11-12: bloqueo
   * tras varios intentos fallidos) y de cuenta activa/desactivada.
   */
  async signIn(email: string, password: string) {
    const { data: lockStatus } = await this.sb.client.rpc('check_account_lock', { p_email: email });
    if (lockStatus?.locked) {
      return { data: { user: null, session: null }, error: { message: 'Cuenta bloqueada temporalmente por varios intentos fallidos. Intenta de nuevo en unos minutos.' } as any };
    }

    const result = await this.sb.client.auth.signInWithPassword({ email, password });

    if (result.error) {
      await this.sb.client.rpc('record_failed_login', { p_email: email });
      return result;
    }

    const userId = result.data.user?.id;
    if (userId) {
      const { data: profile } = await this.sb.client.from('users').select('active').eq('id', userId).single();
      if (profile && profile.active === false) {
        await this.sb.client.auth.signOut();
        return { data: { user: null, session: null }, error: { message: 'Esta cuenta está desactivada. Contacta a un administrador.' } as any };
      }
      await this.sb.client.rpc('record_successful_login', { p_user_id: userId });
    }

    return result;
  }

  /**
   * Challenge & Verify existing MFA code for login
   */
  async verifyMfaLogin(factorId: string, code: string) {
    return this.sb.client.auth.mfa.challengeAndVerify({ factorId, code });
  }

  /**
   * MFA (TOTP) Enrollment Flow
   */
  async getMfaFactors() {
    return this.sb.client.auth.mfa.listFactors();
  }

  async enrollMfa() {
    return this.sb.client.auth.mfa.enroll({ factorType: 'totp' });
  }

  async verifyMfaSetup(factorId: string, code: string) {
    return this.sb.client.auth.mfa.challengeAndVerify({ factorId, code });
  }

  async unenrollMfa(factorId: string) {
    return this.sb.client.auth.mfa.unenroll({ factorId });
  }

  async signOut() {
    const userId = await this.getCurrentUserId();
    if (userId) {
      try {
        await this.sb.client.rpc('record_logout', { p_user_id: userId });
      } catch { /* best-effort logging */ }
    }
    return this.sb.client.auth.signOut();
  }

  /**
   * Cambio de contraseña propio (usuario autenticado). Valida la política
   * de contraseña segura antes de enviarla a Supabase Auth.
   */
  async changePassword(newPassword: string) {
    if (!STRONG_PASSWORD_PATTERN.test(newPassword)) {
      return { data: null, error: { message: STRONG_PASSWORD_HINT } as any };
    }

    const result = await this.sb.client.auth.updateUser({ password: newPassword });
    if (!result.error) {
      const userId = await this.getCurrentUserId();
      if (userId) {
        try {
          await this.sb.client.rpc('record_password_change', { p_user_id: userId });
        } catch { /* best-effort logging */ }
      }
    }
    return result;
  }

  async getSession() {
    const { data, error } = await this.sb.client.auth.getSession();
    return { data, error };
  }

  /**
   * Get current user's ID from the Supabase session.
   */
  async getCurrentUserId(): Promise<string | null> {
    const { data } = await this.sb.client.auth.getSession();
    return data?.session?.user?.id || null;
  }
}
