import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ViewTeams } from './view-teams';

describe('ViewTeams', () => {
  let component: ViewTeams;
  let fixture: ComponentFixture<ViewTeams>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ViewTeams]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ViewTeams);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
