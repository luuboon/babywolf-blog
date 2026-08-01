import { Component, OnInit, inject, ChangeDetectorRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { SupabaseService } from '../../../core/services/supabase.service';
import { AdminApiService } from '../../../core/services/admin-api.service';
import { STRONG_PASSWORD_PATTERN, STRONG_PASSWORD_HINT } from '../../../core/services/auth.service';

interface User {
  id: string;
  email: string;
  username: string;
  role: string;
  active: boolean;
  created_at: string;
}

@Component({
  selector: 'app-user-management',
  standalone: true,
  imports: [CommonModule, FormsModule, RouterModule],
  templateUrl: './user-management.component.html',
  styleUrls: ['./user-management.component.scss']
})
export class UserManagementComponent implements OnInit {
  users: User[] = [];
  loading = true;
  errorMsg = '';

  // Alta de usuario (Practica 11-12, Parte 1)
  showCreateForm = false;
  newEmail = '';
  newUsername = '';
  newPassword = '';
  newRole = 'user';
  creating = false;
  createError = '';
  passwordHint = STRONG_PASSWORD_HINT;

  // Restablecer contraseña por admin
  resettingUserId = '';
  resetPasswordValue = '';
  resetError = '';

  // Edición de usuario por admin (Practica 11-12, Parte 1)
  editingUserId = '';
  editUsernameValue = '';
  editError = '';

  private sb = inject(SupabaseService);
  private adminApi = inject(AdminApiService);
  private cdr = inject(ChangeDetectorRef);

  ngOnInit() {
    this.loadUsers();
  }

  async loadUsers() {
    this.loading = true;
    const { data, error } = await this.sb.client
      .from('users')
      .select('*')
      .order('created_at', { ascending: false });

    if (error) {
      this.errorMsg = 'Error al cargar usuarios: ' + error.message;
      console.error(error);
    } else {
      this.users = data as User[];
    }
    this.loading = false;
    this.cdr.markForCheck();
  }

  async createUser() {
    this.createError = '';
    if (!this.newEmail || !this.newUsername) {
      this.createError = 'Email y nombre de usuario son obligatorios.';
      return;
    }
    if (!STRONG_PASSWORD_PATTERN.test(this.newPassword)) {
      this.createError = this.passwordHint;
      return;
    }

    this.creating = true;
    this.cdr.markForCheck();
    try {
      await this.adminApi.createUser(this.newEmail, this.newPassword, this.newUsername, this.newRole);
      this.showCreateForm = false;
      this.newEmail = this.newUsername = this.newPassword = '';
      this.newRole = 'user';
      await this.loadUsers();
    } catch (err: any) {
      this.createError = err.message || 'Error al crear usuario';
    }
    this.creating = false;
    this.cdr.markForCheck();
  }

  async toggleActive(user: User) {
    const nextActive = !user.active;
    const verb = nextActive ? 'activar' : 'desactivar';
    if (!confirm(`¿Seguro que quieres ${verb} a ${user.email}? (eliminación lógica, no se pierde información)`)) return;

    try {
      await this.adminApi.setActive(user.id, nextActive);
      user.active = nextActive;
    } catch (err: any) {
      alert('Error al actualizar estado: ' + err.message);
    }
    this.cdr.markForCheck();
  }

  async toggleRole(user: User) {
    const newRole = user.role === 'admin' ? 'user' : 'admin';
    try {
      await this.adminApi.updateRole(user.id, newRole);
      user.role = newRole;
    } catch (err: any) {
      alert('Error al actualizar rol: ' + err.message);
    }
    this.cdr.markForCheck();
  }

  startEditUser(user: User) {
    this.editingUserId = user.id;
    this.editUsernameValue = user.username;
    this.editError = '';
  }

  cancelEditUser() {
    this.editingUserId = '';
  }

  async confirmEditUser(user: User) {
    this.editError = '';
    if (!this.editUsernameValue.trim()) {
      this.editError = 'El nombre de usuario no puede estar vacío.';
      return;
    }
    try {
      await this.adminApi.updateUsername(user.id, this.editUsernameValue.trim());
      user.username = this.editUsernameValue.trim();
      this.editingUserId = '';
    } catch (err: any) {
      this.editError = err.message || 'Error al editar el usuario';
    }
    this.cdr.markForCheck();
  }

  startResetPassword(user: User) {
    this.resettingUserId = user.id;
    this.resetPasswordValue = '';
    this.resetError = '';
  }

  cancelResetPassword() {
    this.resettingUserId = '';
  }

  async confirmResetPassword(user: User) {
    this.resetError = '';
    if (!STRONG_PASSWORD_PATTERN.test(this.resetPasswordValue)) {
      this.resetError = this.passwordHint;
      return;
    }
    try {
      await this.adminApi.resetPassword(user.id, this.resetPasswordValue);
      this.resettingUserId = '';
      alert(`Contraseña de ${user.email} restablecida correctamente.`);
    } catch (err: any) {
      this.resetError = err.message || 'Error al restablecer la contraseña';
    }
    this.cdr.markForCheck();
  }
}
