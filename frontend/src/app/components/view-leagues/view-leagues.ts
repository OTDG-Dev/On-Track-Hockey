import { CommonModule } from '@angular/common';
import { Component, signal, WritableSignal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { LeagueData } from '../../interfaces/league-data';
import { LeagueService } from '../../services/league-service';
import { DivisionData } from '../../interfaces/division-data';

@Component({
  selector: 'app-view-leagues',
  imports: [CommonModule, RouterModule, FormsModule],
  templateUrl: './view-leagues.html',
  styleUrl: './view-leagues.css',
})
export class ViewLeagues {

  leagues: WritableSignal<LeagueData[]> = signal([]);
  divisionCounts = signal<Record<number, number>>({});

  constructor(private leagueService: LeagueService) { }

  ngOnInit() {
    this.leagueService.getLeagues().subscribe({
      next: (responseData) => {
        this.leagues.set(responseData.leagues);
        responseData.leagues.forEach(league => {
          this.loadDivisionCount(league.id);
        });
      },
      error: (err) => {
        console.log(err);
      }
    });

  }

  loadDivisionCount(id: number) {
    this.leagueService.getDivisionByLeague(id).subscribe({
      next: (res) => {
        this.divisionCounts.update(counts => ({
          ...counts,
          [id]: res.divisions.length
        }));
      }
    });
  }
}
