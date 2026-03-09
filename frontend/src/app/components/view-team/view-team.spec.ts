import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ViewTeam } from './view-team';

describe('ViewTeam', () => {
  let component: ViewTeam;
  let fixture: ComponentFixture<ViewTeam>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ViewTeam]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ViewTeam);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
