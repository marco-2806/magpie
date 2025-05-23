import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ScraperComponent } from './scraper.component';

describe('ScraperComponent', () => {
  let component: ScraperComponent;
  let fixture: ComponentFixture<ScraperComponent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ScraperComponent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ScraperComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
