import { Component } from '@angular/core';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';

@Component({
  selector: 'app-admin-dashboard',
  standalone: true,
  imports: [CommonModule, RouterModule],
  templateUrl: './admin-dashboard.component.html',
  styleUrls: ['./admin-dashboard.component.scss']
})
export class AdminDashboardComponent {
  
  // Fake stats for UI presentation
  stats = [
    { title: 'Usuarios Totales', value: 124 },
    { title: 'Posts Publicados', value: 45 },
    { title: 'Comentarios', value: 302 },
    { title: 'Reportes Activos', value: 2 }
  ];

}
