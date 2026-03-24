import { Component, signal, WritableSignal } from '@angular/core';
import { LeagueService } from '../../services/league-service';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';

@Component({
  selector: 'app-edit-league',
  imports: [FormsModule],
  templateUrl: './edit-league.html',
  styleUrl: './edit-league.css',
})
export class EditLeague {

  league_id: number = -1
  name: WritableSignal<string> = signal("");

  successMessage: WritableSignal<string> = signal('');
  errorMessage: WritableSignal<string> = signal('');
  isFading = signal(false);

  constructor(private leagueService: LeagueService, private route: ActivatedRoute, private router: Router) {}

  ngOnInit()
  {
    const id = this.route.snapshot.paramMap.get('id');
    this.league_id = Number(id);

    this.getLeague(this.league_id);
  }

  getLeague(id: number) {
    this.leagueService.getLeague(id)
    .subscribe(
      {
        next: (responseData) => {
          this.name.set(responseData.league.name);
        },
        error: (err) => {
          console.log(err);
          this.router.navigate(['/view-divisions']);
        }
      }
    )
  }

  patchLeague() {
    this.leagueService.patchLeague(this.name(), this.league_id)
    .subscribe(
      {
        next: (responseData) => {
          this.successMessage.set(`Edited ${responseData.league.name} League`);

          setTimeout(() => {
            this.isFading.set(true);
          }, 2500);

          setTimeout(() => {
            this.successMessage.set('');
            this.isFading.set(false);
          }, 2750);
        },
        error: (err) => {
          this.errorMessage.set(`Failed to Edit League`);

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
