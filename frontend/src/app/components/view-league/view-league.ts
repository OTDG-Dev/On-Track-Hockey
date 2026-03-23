import { CommonModule } from '@angular/common';
import { Component, signal, WritableSignal } from '@angular/core';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { DivisionData } from '../../interfaces/division-data';
import { LeagueService } from '../../services/league-service';
import { DivisionService } from '../../services/division-service';
import { LeagueData } from '../../interfaces/league-data';

@Component({
  selector: 'app-view-league',
  imports: [RouterModule, CommonModule],
  templateUrl: './view-league.html',
  styleUrl: './view-league.css',
})
export class ViewLeague {

  activeTab: 'info' | 'divisions' = 'info';
  leagueId: number = -1;
  name: WritableSignal<string> = signal("");
  divisions: WritableSignal<DivisionData[]> = signal([]);

  avatarUrl: WritableSignal<string> = signal("https://a.espncdn.com/combiner/i?img=/i/headshots/nhl/players/full/5149125.png&w=350&h=254");

  constructor(private leagueService: LeagueService, private divisionService: DivisionService, private route: ActivatedRoute, private router: Router) {}

  ngOnInit(){
    const id = this.route.snapshot.paramMap.get('id');
    this.leagueId = Number(id);

    this.getLeague(this.leagueId);
  }

  getLeague(id: number) {
    this.leagueService.getLeague(id)
    .subscribe(
      {
        next: (responseData) => {
          this.name.set(responseData.league.name);
        }
      }
    )
  }

}
