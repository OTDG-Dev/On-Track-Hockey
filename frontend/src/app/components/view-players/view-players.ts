import { Component, signal, WritableSignal } from '@angular/core';
import { PlayerService } from '../../services/player-service';
import { PlayerData } from '../../interfaces/player-data';
import { CommonModule } from '@angular/common';
import { RouterLink, RouterModule } from "@angular/router";
import { TeamService } from '../../services/team-service';
import { TeamData } from '../../interfaces/team-data';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-view-players',
  standalone: true,
  imports: [CommonModule, RouterModule, FormsModule],
  templateUrl: './view-players.html',
  styleUrl: './view-players.css',
})
export class ViewPlayers {

  players: WritableSignal<PlayerData[]> = signal([]);
  teams: WritableSignal<TeamData[]> = signal([]);

  selectedPosition: string = '';
  current_team_id: number | null = null;

  constructor(private playerService: PlayerService, private teamService: TeamService) { }

  ngOnInit() {
    this.playerService.getPlayers().subscribe({
      next: (responseData) => {
        this.players.set(responseData.players);
      },
      error: (err) => {
        console.log(err);
      }
    });

    this.teamService.getTeams().subscribe({
      next: (responseData) => {
        this.teams.set(responseData.teams);
      },
      error: (err) => {
        console.log(err);
      }
    });
  }

  applyFilters() {
    this.playerService
      .getPlayers(this.selectedPosition, this.current_team_id)
      .subscribe({
        next: (responseData) => {
          this.players.set(responseData.players);
        },
        error: (err) => {
          console.error(err);
        }
      });
  }

}
