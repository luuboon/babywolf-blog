import { RenderMode, ServerRoute } from '@angular/ssr';

export const serverRoutes: ServerRoute[] = [
  {
    path: 'posts/:slug',
    renderMode: RenderMode.Server
  },
  {
    path: 'category/:name',
    renderMode: RenderMode.Server
  },
  {
    path: 'editor/:id',
    renderMode: RenderMode.Server
  },
  {
    path: 'search',
    renderMode: RenderMode.Server
  },
  {
    path: '**',
    renderMode: RenderMode.Prerender
  }
];
