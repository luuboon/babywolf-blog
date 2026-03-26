import { Routes } from '@angular/router';
import { AuthGuard } from './core/guards/auth.guard';

export const routes: Routes = [
  // ===== AUTH (standalone, no layout wrapper) =====
  {
    path: 'login',
    loadComponent: () => import('./features/auth/login.component').then(m => m.LoginComponent)
  },
  {
    path: 'register',
    loadComponent: () => import('./features/auth/register.component').then(m => m.RegisterComponent)
  },

  // ===== ADMIN (AdminLayout wrapper) =====
  {
    path: 'admin',
    loadComponent: () => import('./features/admin/layout/admin-layout').then(m => m.AdminLayoutComponent),
    canActivate: [AuthGuard],
    children: [
      {
        path: '',
        loadComponent: () => import('./features/admin/admin-dashboard.component').then(m => m.AdminDashboardComponent)
      },
      {
        path: 'users',
        loadComponent: () => import('./features/admin/users/user-management.component').then(m => m.UserManagementComponent)
      },
      {
        path: 'posts',
        loadComponent: () => import('./features/admin/posts/post-management.component').then(m => m.PostManagementComponent)
      },
      {
        path: 'comments',
        loadComponent: () => import('./features/admin/comments/comment-management.component').then(m => m.CommentManagementComponent)
      },
      {
        path: 'editor/new',
        loadComponent: () => import('./features/admin/posts/editor.component').then(m => m.PostEditorComponent)
      }
    ]
  },

  // ===== PUBLIC (UserLayout wrapper) =====
  {
    path: '',
    loadComponent: () => import('./features/user/layout/user-layout').then(m => m.UserLayoutComponent),
    children: [
      {
        path: '',
        loadComponent: () => import('./features/home/presentation/pages/home.component').then(m => m.HomeComponent)
      },
      {
        path: 'posts',
        loadComponent: () => import('./features/posts/presentation/pages/posts-page/posts-page').then(m => m.PostsPage)
      },
      {
        path: 'posts/:slug',
        loadComponent: () => import('./features/posts/presentation/pages/post-detail/post-detail').then(m => m.PostDetailPage)
      },
      {
        path: 'category/:name',
        loadComponent: () => import('./features/category/category-posts').then(m => m.CategoryPostsPage)
      },
      {
        path: 'search',
        loadComponent: () => import('./features/search/search-results').then(m => m.SearchResultsPage)
      },
      {
        path: 'terms',
        loadComponent: () => import('./features/legal/presentation/pages/terms-page/terms-page').then(m => m.TermsPage)
      },
      {
        path: 'profile',
        loadComponent: () => import('./features/user/profile/profile.page').then(m => m.ProfilePage),
        canActivate: [AuthGuard]
      },
      {
        path: 'editor/new',
        loadComponent: () => import('./features/user/editor/user-editor').then(m => m.UserEditorComponent),
        canActivate: [AuthGuard]
      },
      {
        path: 'editor/:id',
        loadComponent: () => import('./features/user/editor/user-editor').then(m => m.UserEditorComponent),
        canActivate: [AuthGuard]
      }
    ]
  }
];
