import { Component, signal, WritableSignal } from '@angular/core';
import { ActivatedRoute, Router, RouterModule } from '@angular/router';
import { TeamService } from '../../services/team-service';
import { DivisionData } from '../../interfaces/division-data';
import { DivisionService } from '../../services/division-service';
import { CommonModule } from '@angular/common';
import { RosterService } from '../../services/roster-service';
import { PlayerData } from '../../interfaces/player-data';

@Component({
  selector: 'app-view-team',
  imports: [RouterModule, CommonModule],
  templateUrl: './view-team.html',
  styleUrl: './view-team.css',
})
export class ViewTeam {

  activeTab: 'info' | 'players' = 'info';

  teamId: number = -1;
  full_name: WritableSignal<string> = signal("");
  short_name: WritableSignal<string> = signal("");
  division_id: WritableSignal<number> = signal(-1);
  is_active: WritableSignal<boolean> = signal(false);
  divisions: WritableSignal<DivisionData[]> = signal([]);

  forwards: WritableSignal<PlayerData[]> = signal([]);
  defenseman: WritableSignal<PlayerData[]> = signal([]);
  goalies: WritableSignal<PlayerData[]> = signal([]);

  avatarUrl: WritableSignal<string> = signal("https://a.espncdn.com/combiner/i?img=/i/headshots/nhl/players/full/5149125.png&w=350&h=254");

  constructor(private teamService: TeamService, private divisionService: DivisionService, private rosterService: RosterService, private route: ActivatedRoute, private router: Router) {}

  ngOnInit(){
    const id = this.route.snapshot.paramMap.get('id');
    this.teamId = Number(id);

    this.getTeam(this.teamId);

    this.divisionService.getDivisions().subscribe({
      next: (responseData) => {
        this.divisions.set(responseData.divisions);
      },
      error: (err) => {
        console.log(err);
      }
    });

    this.rosterService.getRoster(this.teamId).subscribe({
      next: (responseData) => {
        this.forwards.set(responseData.roster.forwards);
        this.defenseman.set(responseData.roster.defensemen);
        this.goalies.set(responseData.roster.goalies);
        console.log(this.forwards());
        console.log(this.defenseman());
        console.log(this.goalies());
      },
      error: (err) => {
        console.log(err);
      }
    })
  }

  getTeam(id: number) {
    this.teamService.getTeam(id)
    .subscribe(
      {
        next: (responseData) => {
          this.full_name.set(responseData.team.full_name);
          this.short_name.set(responseData.team.short_name);
          this.is_active.set(responseData.team.is_active);
          this.division_id.set(responseData.team.division_id);
          console.log("is_active : " + this.is_active());
          console.log(responseData.team);
        },
        error: (err) => {
          console.log(err);
          this.router.navigate(['/view-teams']);
        }
      }
    )
  }

  getDivisionName(id: number) {
    const division = this.divisions().find(d => d.id === id);
    return division ? division.name : 'Unknown';
  }

  deactivateTeam() {
    this.teamService.deactivateTeam(this.teamId).subscribe({
      next: () => {
        this.router.navigate(['/view-teams']);
      },
      error: (err) => {
        console.log(err);
      }
    });
  }

}
