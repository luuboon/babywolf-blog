import { Routes } from '@angular/router';
import { AuthGuard } from './core/guards/auth.guard';


export const routes: Routes = [
  {
    path: '',
    loadComponent: () => import('./features/home/presentation/pages/home.component').then(m => m.HomeComponent)
  },
  {
    path: 'login',
    loadComponent: () => import('./features/auth/login.component').then(m => m.LoginComponent)
  },
  {
    path: 'register',
    loadComponent: () => import('./features/auth/register.component').then(m => m.RegisterComponent)
  },
  {
    path: 'admin',
    loadComponent: () => import('./features/admin/admin-dashboard.component').then(m => m.AdminDashboardComponent),
    canActivate: [AuthGuard]
  },
  {
    path: 'admin/users',
    loadComponent: () => import('./features/admin/users/user-management.component').then(m => m.UserManagementComponent),
    canActivate: [AuthGuard]
  },
  {
    path: 'admin/posts',
    loadComponent: () => import('./features/admin/posts/post-management.component').then(m => m.PostManagementComponent),
    canActivate: [AuthGuard]
  },
  {
    path: 'admin/comments',
    loadComponent: () => import('./features/admin/comments/comment-management.component').then(m => m.CommentManagementComponent),
    canActivate: [AuthGuard]
  },
  {
    path: 'editor/new',
    loadComponent: () => import('./features/admin/posts/editor.component').then(m => m.PostEditorComponent),
    canActivate: [AuthGuard]
  },
  {
    path: 'posts',
    loadComponent: () => import('./features/posts/presentation/pages/posts-page/posts-page').then(m => m.PostsPage)
  },
  {
    path: 'contact',
    loadComponent: () => import('./features/home/presentation/pages/home.component').then(m => m.HomeComponent) // Temp placeholder
  },
  {
    path: 'terms',
    loadComponent: () => import('./features/legal/presentation/pages/terms-page/terms-page').then(m => m.TermsPage)
  }
];
