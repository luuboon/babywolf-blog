import { Component } from '@angular/core';
import { NgFor } from '@angular/common';

interface ExternalLink {
  label: string;
  url: string;
}

@Component({
  selector: 'app-external-links',
  standalone: true,
  imports: [NgFor],
  templateUrl: './external-links.html',
  styleUrls: ['./external-links.scss'],
})

export class ExternalLinks {
  links: ExternalLink[] = [
    { label: 'GitHub', url: 'https://github.com/luuboon/babywolf-blog' },
    { label: 'LinkedIn', url: 'https://www.linkedin.com/in/abraham-duran-50828139b/' },
    { label: 'Docs', url: 'https://angular.dev' },
  ];
}
