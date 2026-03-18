import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ViewLeagues } from './view-leagues';

describe('ViewLeagues', () => {
  let component: ViewLeagues;
  let fixture: ComponentFixture<ViewLeagues>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ViewLeagues]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ViewLeagues);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
