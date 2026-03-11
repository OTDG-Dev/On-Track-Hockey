import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ViewDivision } from './view-division';

describe('ViewDivision', () => {
  let component: ViewDivision;
  let fixture: ComponentFixture<ViewDivision>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ViewDivision]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ViewDivision);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
