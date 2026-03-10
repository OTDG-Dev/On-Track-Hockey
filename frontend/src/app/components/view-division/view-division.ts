import { Component, signal, WritableSignal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { LeagueData } from '../../interfaces/league-data';
import { DivisionService } from '../../services/division-service';
import { LeagueService } from '../../services/league-service';
import { ActivatedRoute, Router } from '@angular/router';

@Component({
  selector: 'app-view-division',
  imports: [FormsModule],
  templateUrl: './view-division.html',
  styleUrl: './view-division.css',
})
export class ViewDivision {

  division_id: number = 1;
  name: WritableSignal<string> = signal("");
  league_id: WritableSignal<number> = signal(-1);

  leagues: WritableSignal<LeagueData[]> = signal([]);

  successMessage: WritableSignal<string> = signal('');
  errorMessage: WritableSignal<string> = signal('');
  isFading = signal(false);

  constructor(private divisionService: DivisionService, private leageService: LeagueService, private route: ActivatedRoute, private router: Router) {}

  ngOnInit()
  {
    const id = this.route.snapshot.paramMap.get('id');
    this.division_id = Number(id);

    this.getDivision(this.division_id);

    this.leageService.getLeagues()
    .subscribe(
      {
        next: (responseData) => {
          this.leagues.set(responseData.leagues);
        },
        error: (err) => {
          console.log(err);
        }
      }
    )
  }

  getDivision(id: number) {
    this.divisionService.getDivision(id)
    .subscribe(
      {
        next: (responseData) => {
          this.name.set(responseData.division.name);
          this.league_id.set(responseData.division.league_id);
        },
        error: (err) => {
          console.log(err);
          this.router.navigate(['/view-divisions']);
        }
      }
    )
  }

  patchDivision() 
  {
    this.divisionService.patchDivision(this.league_id(), this.name(), this.division_id)
    .subscribe(
      {
        next: (responseData) => {
          this.successMessage.set(`Edited ${responseData.division.name} Division`);

          setTimeout(() => {
            this.isFading.set(true);
          }, 2500);

          setTimeout(() => {
            this.successMessage.set('');
            this.isFading.set(false);
          }, 2750);
        },
        error: (err) => {
          this.errorMessage.set(`Failed to Edit Division`);

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
