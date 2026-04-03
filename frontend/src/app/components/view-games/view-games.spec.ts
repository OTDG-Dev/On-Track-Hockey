import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ViewGames } from './view-games';

describe('ViewGames', () => {
  let component: ViewGames;
  let fixture: ComponentFixture<ViewGames>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ViewGames]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ViewGames);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
