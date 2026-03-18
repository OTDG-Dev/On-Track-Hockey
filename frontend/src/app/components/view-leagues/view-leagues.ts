import { CommonModule } from '@angular/common';
import { Component, signal, WritableSignal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { LeagueData } from '../../interfaces/league-data';
import { LeagueService } from '../../services/league-service';

@Component({
  selector: 'app-view-leagues',
  imports: [CommonModule, RouterModule, FormsModule],
  templateUrl: './view-leagues.html',
  styleUrl: './view-leagues.css',
})
export class ViewLeagues {

  leagues: WritableSignal<LeagueData[]> = signal([]);

  constructor(private leagueService: LeagueService) { }

  ngOnInit() {
    this.leagueService.getLeagues().subscribe({
      next: (responseData) => {
        this.leagues.set(responseData.leagues);
      },
      error: (err) => {
        console.log(err);
      }
    });
  }

  getNoDivisions(id: number) {
    return id;
  }
}
