import { ComponentFixture, TestBed } from '@angular/core/testing';

import { CreatePlayer } from './create-player';

describe('CreatePlayer', () => {
  let component: CreatePlayer;
  let fixture: ComponentFixture<CreatePlayer>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [CreatePlayer]
    })
    .compileComponents();

    fixture = TestBed.createComponent(CreatePlayer);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
