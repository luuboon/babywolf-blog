import { Routes } from '@angular/router';
import { Component } from '@angular/core';
import { HomeComponent } from './features/home/presentation/pages/home.component';

@Component({ template: `<h1>Posts</h1>` })
class PostsPage { }

@Component({ template: `<h1>Contact</h1>` })
class ContactPage { }

export const routes: Routes = [
  { path: '', component: HomeComponent },
  { path: 'posts', component: PostsPage },
  { path: 'contact', component: ContactPage },
];
