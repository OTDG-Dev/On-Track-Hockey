import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CreateDivision } from './create-division';

describe('CreateDivision', () => {
  let component: CreateDivision;
  let fixture: ComponentFixture<CreateDivision>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CreateDivision]
    })
    .compileComponents();

    fixture = TestBed.createComponent(CreateDivision);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
