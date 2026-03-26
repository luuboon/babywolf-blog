import { Component, OnInit, inject, ChangeDetectorRef } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { SupabaseService } from '../../../core/services/supabase.service';

interface User {
  id: string;
  email: string;
  username: string;
  role: string;
  created_at: string;
}

@Component({
  selector: 'app-user-management',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './user-management.component.html',
  styleUrls: ['./user-management.component.scss']
})
export class UserManagementComponent implements OnInit {
  users: User[] = [];
  loading = true;
  errorMsg = '';

  private sb = inject(SupabaseService);
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

  async deleteUser(id: string) {
    if (!confirm('¿Estás seguro de eliminar este usuario?')) return;
    
    const { error } = await this.sb.client
      .from('users')
      .delete()
      .eq('id', id);

    if (error) {
      alert('Error al eliminar usuario: ' + error.message);
      console.error(error);
    } else {
      this.users = this.users.filter(u => u.id !== id);
    }
    this.cdr.markForCheck();
  }

  async toggleRole(user: User) {
    const newRole = user.role === 'admin' ? 'user' : 'admin';
    
    const { error } = await this.sb.client
      .from('users')
      .update({ role: newRole })
      .eq('id', user.id);

    if (error) {
      alert('Error al actualizar rol: ' + error.message);
      console.error(error);
    } else {
      user.role = newRole;
    }
    this.cdr.markForCheck();
  }
}
