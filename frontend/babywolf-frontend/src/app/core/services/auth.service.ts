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

  constructor() {
    // Check active session on init
    this.sb.client.auth.getSession().then(({ data }) => {
       if (data.session) {
         this.currentUserSubject.next(data.session.user);
       }
    });

    this.sb.client.auth.onAuthStateChange((_event, session) => {
      if (session) {
        this.currentUserSubject.next(session.user);
      } else {
        this.currentUserSubject.next(null);
      }
    });
  }

  async signUp(email: string, password: string) {
    return this.sb.client.auth.signUp({ email, password });
  }

  async signIn(email: string, password: string) {
    return this.sb.client.auth.signInWithPassword({ email, password });
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
