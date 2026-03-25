import { Injectable, inject } from '@angular/core';
import { BehaviorSubject, Observable } from 'rxjs';
import { User } from '@supabase/supabase-js';
import { SupabaseService } from './supabase.service';

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

  async signIn(email: string, password: string) {
    return this.sb.client.auth.signInWithPassword({ email, password });
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
    return this.sb.client.auth.signOut();
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
