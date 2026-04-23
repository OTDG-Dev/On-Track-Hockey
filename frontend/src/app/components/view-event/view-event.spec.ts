import { ComponentFixture, TestBed } from '@angular/core/testing';

import { ViewEvent } from './view-event';

describe('ViewEvent', () => {
  let component: ViewEvent;
  let fixture: ComponentFixture<ViewEvent>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [ViewEvent]
    })
    .compileComponents();

    fixture = TestBed.createComponent(ViewEvent);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
