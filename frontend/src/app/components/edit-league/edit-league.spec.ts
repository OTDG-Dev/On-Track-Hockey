import { ComponentFixture, TestBed } from '@angular/core/testing';

import { EditLeague } from './edit-league';

describe('EditLeague', () => {
  let component: EditLeague;
  let fixture: ComponentFixture<EditLeague>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [EditLeague]
    })
    .compileComponents();

    fixture = TestBed.createComponent(EditLeague);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
