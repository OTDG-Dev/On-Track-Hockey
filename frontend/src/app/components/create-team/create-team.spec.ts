import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CreateTeam } from './create-team';

describe('CreateTeam', () => {
  let component: CreateTeam;
  let fixture: ComponentFixture<CreateTeam>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CreateTeam]
    })
    .compileComponents();

    fixture = TestBed.createComponent(CreateTeam);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
