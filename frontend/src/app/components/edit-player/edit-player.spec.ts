import { ComponentFixture, TestBed } from '@angular/core/testing';

import { EditPlayer } from './edit-player';

describe('EditPlayer', () => {
  let component: EditPlayer;
  let fixture: ComponentFixture<EditPlayer>;

  beforeEach(async () => {
    await TestBed.configureTestingModule({
      imports: [EditPlayer]
    })
    .compileComponents();

    fixture = TestBed.createComponent(EditPlayer);
    component = fixture.componentInstance;
    await fixture.whenStable();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
