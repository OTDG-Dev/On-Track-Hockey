import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CreateLeague } from './create-league';

describe('CreateLeague', () => {
  let component: CreateLeague;
  let fixture: ComponentFixture<CreateLeague>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CreateLeague]
    })
    .compileComponents();

    fixture = TestBed.createComponent(CreateLeague);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
