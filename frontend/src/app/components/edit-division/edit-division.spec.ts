import { ComponentFixture, TestBed } from '@angular/core/testing';

import { EditDivision } from './edit-division';

describe('ViewDivision', () => {
  let component: EditDivision;
  let fixture: ComponentFixture<EditDivision>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [EditDivision]
    })
    .compileComponents();

    fixture = TestBed.createComponent(EditDivision);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
