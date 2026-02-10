import { Component } from '@angular/core';

@Component({
    selector: 'app-home',
    standalone: true,
    template: `
    <div class="home-container">
      <h1>BabyWolf Blog</h1>
      <p>Your space to share your thoughts</p>
    </div>
  `,
    styles: [`
    .home-container {
      padding: 2rem;
      text-align: center;
    }
  `]
})
export class HomeComponent { }
