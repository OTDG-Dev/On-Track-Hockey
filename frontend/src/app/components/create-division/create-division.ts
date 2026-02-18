import { Component, signal, WritableSignal } from '@angular/core';
import { DivisionService } from '../../services/division-service';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-create-division',
  imports: [FormsModule],
  templateUrl: './create-division.html',
  styleUrl: './create-division.css',
})
export class CreateDivision {

  league_id: number = 1;
  name: string = "";

  successMessage: WritableSignal<string> = signal('');
  errorMessage: WritableSignal<string> = signal('');
  isFading = signal(false);

  constructor(private divisionService: DivisionService) {}

  postDivision() 
  {
    this.divisionService.createDivision(this.league_id, this.name)
    .subscribe(
      {
        next: (responseData) => {
          this.successMessage.set(`Created ${responseData.division.name} Division`);

          setTimeout(() => {
            this.isFading.set(true);
          }, 2500);

          setTimeout(() => {
            this.successMessage.set('');
            this.isFading.set(false);
          }, 2750);
        },
        error: (err) => {
          this.errorMessage.set(`Failed to Create Division`)

          setTimeout(() => {
            this.isFading.set(true);
          }, 2500);

          setTimeout(() => {
            this.errorMessage.set('');
            this.isFading.set(false);
          }, 2750);
        }
      }
    )
  }



}
