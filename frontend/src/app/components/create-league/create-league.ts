import { Component, signal, WritableSignal } from '@angular/core';
import { LeagueService } from '../../services/league-service';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-create-league',
  imports: [FormsModule],
  templateUrl: './create-league.html',
  styleUrl: './create-league.css',
})
export class CreateLeague {

  name: string = "";

  successMessage: WritableSignal<string> = signal('');
  errorMessage: WritableSignal<string> = signal('');
  isFading = signal(false);

  constructor(private leagueService: LeagueService) {}

  postLeague()
  {
    this.leagueService.createLeague(this.name)
    .subscribe
    ({
      next: (responseData) => {
        this.successMessage.set(`Created ${responseData.league.name} League`);

        setTimeout(() => {
          this.isFading.set(true);
        }, 2500);

        setTimeout(() => {
          this.successMessage.set('');
          this.isFading.set(false);
        }, 2750);
      },
      error: (err) => {
        this.errorMessage.set(`Failed to Create League`);

        setTimeout(() => {
          this.isFading.set(true);
        }, 2500);

        setTimeout(() => {
          this.errorMessage.set('');
          this.isFading.set(false);
        }, 2750);
      }
    })
  }
}
