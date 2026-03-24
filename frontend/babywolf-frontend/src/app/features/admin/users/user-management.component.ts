import { Component, OnInit } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { HttpClient } from '@angular/common/http';
import { environment } from '../../../../environments/environment';

interface User {
  id: string;
  email: string;
  role: string;
  created_at: string;
}

@Component({
  selector: 'app-user-management',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './user-management.component.html',
  styleUrls: ['../admin-dashboard.component.scss']
})
export class UserManagementComponent implements OnInit {
  users: User[] = [];
  loading = true;
  errorMsg = '';

  constructor(private http: HttpClient) {}

  ngOnInit() {
    this.loadUsers();
  }

  loadUsers() {
    this.loading = true;
    this.http.get<User[]>(`${environment.apiUrl}/admin/users`).subscribe({
      next: (data) => {
        this.users = data;
        this.loading = false;
      },
      error: (err) => {
        this.errorMsg = 'Error al cargar usuarios. ¿Eres administrador?';
        this.loading = false;
        console.error(err);
      }
    });
  }

  deleteUser(id: string) {
    if (!confirm('¿Estás seguro de eliminar este usuario?')) return;
    
    this.http.delete(`${environment.apiUrl}/admin/users/${id}`).subscribe({
      next: () => {
        this.users = this.users.filter(u => u.id !== id);
      },
      error: (err) => {
        alert('Error al eliminar usuario');
        console.error(err);
      }
    });
  }

  toggleRole(user: User) {
    const newRole = user.role === 'admin' ? 'user' : 'admin';
    
    this.http.put(`${environment.apiUrl}/admin/users/${user.id}/role`, { role: newRole }).subscribe({
      next: () => {
        user.role = newRole;
      },
      error: (err) => {
        alert('Error al actualizar rol');
        console.error(err);
      }
    });
  }
}
