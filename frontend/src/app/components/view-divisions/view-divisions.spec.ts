import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ViewDivisions } from './view-divisions';

describe('ViewDivisions', () => {
  let component: ViewDivisions;
  let fixture: ComponentFixture<ViewDivisions>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ViewDivisions]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ViewDivisions);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
