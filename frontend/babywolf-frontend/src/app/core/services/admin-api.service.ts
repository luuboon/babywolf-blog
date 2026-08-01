import { Injectable, inject } from '@angular/core';
import { SupabaseService } from './supabase.service';
import { environment } from '../../../environments/environment';

/**
 * Cliente para los endpoints administrativos del backend Go
 * (Practica 11-12): operaciones que requieren la service_role key y por
 * lo tanto no pueden hacerse directo desde el navegador con Supabase.
 */
@Injectable({ providedIn: 'root' })
export class AdminApiService {
  private sb = inject(SupabaseService);
  private baseUrl = environment.backendUrl;

  private async authHeaders(): Promise<Record<string, string>> {
    const { data } = await this.sb.client.auth.getSession();
    return {
      'Content-Type': 'application/json',
      Authorization: `Bearer ${data.session?.access_token ?? ''}`
    };
  }

  private async request(method: string, path: string, body?: unknown) {
    const res = await fetch(`${this.baseUrl}${path}`, {
      method,
      headers: await this.authHeaders(),
      body: body !== undefined ? JSON.stringify(body) : undefined
    });
    const json = await res.json().catch(() => ({}));
    if (!res.ok) throw new Error(json.error || `Error ${res.status}`);
    return json;
  }

  createUser(email: string, password: string, username: string, role: string) {
    return this.request('POST', '/admin/users', { email, password, username, role });
  }

  setActive(id: string, active: boolean) {
    return this.request('PATCH', `/admin/users/${id}/active`, { active });
  }

  updateRole(id: string, role: string) {
    return this.request('PATCH', `/admin/users/${id}/role`, { role });
  }

  updateUsername(id: string, username: string) {
    return this.request('PATCH', `/admin/users/${id}`, { username });
  }

  resetPassword(id: string, password: string) {
    return this.request('PATCH', `/admin/users/${id}/password`, { password });
  }
}
