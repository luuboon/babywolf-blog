import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ExternalLinks } from './external-links';

describe('ExternalLinks', () => {
  let component: ExternalLinks;
  let fixture: ComponentFixture<ExternalLinks>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ExternalLinks]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ExternalLinks);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
