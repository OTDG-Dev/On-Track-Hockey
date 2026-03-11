import { Component, signal, WritableSignal } from '@angular/core';
import { TeamData } from '../../interfaces/team-data';
import { DivisionData } from '../../interfaces/division-data';
import { TeamService } from '../../services/team-service';
import { DivisionService } from '../../services/division-service';
import { CommonModule } from '@angular/common';
import { RouterModule } from '@angular/router';
import { FormsModule } from '@angular/forms';
import { LeagueData } from '../../interfaces/league-data';
import { LeagueService } from '../../services/league-service';

@Component({
  selector: 'app-view-divisions',
  imports: [CommonModule, RouterModule, FormsModule],
  templateUrl: './view-divisions.html',
  styleUrl: './view-divisions.css',
})
export class ViewDivisions {

  leagues: WritableSignal<LeagueData[]> = signal([]);
  divisions: WritableSignal<DivisionData[]> = signal([]);

  constructor(private leagueService: LeagueService, private divisionService: DivisionService) { }

  ngOnInit() {
    this.leagueService.getLeagues().subscribe({
      next: (responseData) => {
        this.leagues.set(responseData.leagues);
      },
      error: (err) => {
        console.log(err);
      }
    });

    this.divisionService.getDivisions().subscribe({
      next: (responseData) => {
        this.divisions.set(responseData.divisions);
      },
      error: (err) => {
        console.log(err);
      }
    });
  }

  getLeagueName(id: number) {
    const league = this.leagues().find(d => d.id === id);
    return league ? league.name : 'Unknown';
  }

}
