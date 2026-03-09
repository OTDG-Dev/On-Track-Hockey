import { CommonModule } from '@angular/common';
import { Component, signal, WritableSignal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { RouterModule } from '@angular/router';
import { TeamData } from '../../interfaces/team-data';
import { DivisionData } from '../../interfaces/division-data';
import { DivisionService } from '../../services/division-service';
import { TeamService } from '../../services/team-service';

@Component({
  selector: 'app-view-teams',
  imports: [CommonModule, RouterModule, FormsModule],
  templateUrl: './view-teams.html',
  styleUrl: './view-teams.css',
})
export class ViewTeams {

  teams: WritableSignal<TeamData[]> = signal([]);
  divisions: WritableSignal<DivisionData[]> = signal([]);

  constructor(private teamService: TeamService, private divisionService: DivisionService) { }

  ngOnInit() {
    this.teamService.getTeams().subscribe({
      next: (responseData) => {
        this.teams.set(responseData.teams);
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

  getDivisionName(id: number) {
    const division = this.divisions().find(d => d.id === id);
    return division ? division.name : 'Unknown';
  }
}
