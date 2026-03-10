import { ComponentFixture, TestBed } from '@angular/core/testing';

import { EditTeam } from './edit-team';

describe('EditTeam', () => {
  let component: EditTeam;
  let fixture: ComponentFixture<EditTeam>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [EditTeam]
    })
    .compileComponents();

    fixture = TestBed.createComponent(EditTeam);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
