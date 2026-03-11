import { Component, signal, WritableSignal } from '@angular/core';
import { TeamService } from '../../services/team-service';
import { DivisionService } from '../../services/division-service';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { DivisionData } from '../../interfaces/division-data';
import { CommonModule } from '@angular/common';
import { LeagueService } from '../../services/league-service';
import { LeagueData } from '../../interfaces/league-data';

@Component({
  selector: 'app-view-division',
  imports: [RouterModule, CommonModule],
  templateUrl: './view-division.html',
  styleUrl: './view-division.css',
})
export class ViewDivision {

  divisionId: number = -1;
  name: WritableSignal<string> = signal("");
  leagueId: WritableSignal<number> = signal(-1);
  leagues: WritableSignal<LeagueData[]> = signal([]);

  avatarUrl: WritableSignal<string> = signal("https://a.espncdn.com/combiner/i?img=/i/headshots/nhl/players/full/5149125.png&w=350&h=254");

  constructor(private leagueService: LeagueService, private divisionService: DivisionService, private route: ActivatedRoute, private router: Router) {}

  ngOnInit(){
    const id = this.route.snapshot.paramMap.get('id');
    this.divisionId = Number(id);

    this.getDivision(this.divisionId);

    this.leagueService.getLeagues().subscribe({
      next: (responseData) => {
        this.leagues.set(responseData.leagues);
      },
      error: (err) => {
        console.log(err);
      }
    });
  }

  getDivision(id: number) {
    this.divisionService.getDivision(id)
    .subscribe(
      {
        next: (responseData) => {
          this.name.set(responseData.division.name);
          this.leagueId.set(responseData.division.league_id);
          console.log(responseData.division);
        },
        error: (err) => {
          console.log(err);
          this.router.navigate(['/view-divisions']);
        }
      }
    )
  }

  getLeagueName(id: number) {
    const league = this.leagues().find(d => d.id === id);
    return league ? league.name : 'Unknown';
  }

}
