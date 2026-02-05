import { Routes } from '@angular/router';
import { Component } from '@angular/core';

@Component({ template: `<h1>Home</h1>` })
class HomePage {}

@Component({ template: `<h1>Posts</h1>` })
class PostsPage {}

@Component({ template: `<h1>Contact</h1>` })
class ContactPage {}

export const routes: Routes = [
  { path: '', component: HomePage },
  { path: 'posts', component: PostsPage },
  { path: 'contact', component: ContactPage },
];
