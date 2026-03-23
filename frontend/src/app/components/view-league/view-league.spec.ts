import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ViewLeague } from './view-league';

describe('ViewLeague', () => {
  let component: ViewLeague;
  let fixture: ComponentFixture<ViewLeague>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ViewLeague]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ViewLeague);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
