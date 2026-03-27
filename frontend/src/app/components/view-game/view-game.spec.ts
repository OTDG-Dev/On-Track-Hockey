import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ViewGame } from './view-game';

describe('ViewGame', () => {
  let component: ViewGame;
  let fixture: ComponentFixture<ViewGame>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ViewGame]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ViewGame);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
