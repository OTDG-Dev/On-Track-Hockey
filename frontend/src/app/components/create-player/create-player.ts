import { Component, signal, WritableSignal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { PlayerService } from '../../services/player-service';
import { CommonModule } from '@angular/common';
import { TeamService } from '../../services/team-service';
import { TeamData } from '../../interfaces/team-data';

@Component({
  selector: 'app-create-player',
  imports: [FormsModule, CommonModule],
  templateUrl: './create-player.html',
  styleUrl: './create-player.css',
})
export class CreatePlayer {

  firstName: string = "";
  lastName: string = "";
  sweaterNumber: string = "";
  position: string = "";
  handedness: string = "";
  birthCountry: string = "";
  dob: string = "";
  current_team_id: number = 1;

  teams: WritableSignal<TeamData[]> = signal([]);

  successMessage: WritableSignal<string> = signal('');
  errorMessage: WritableSignal<string> = signal('');
  isFading = signal(false);

  constructor(private playerService: PlayerService, private teamService: TeamService) { }

  ngOnInit() {
    this.teamService.getTeams().subscribe({
      next: (responseData) => {
        this.teams.set(responseData.teams);
      },
      error: (err) => {
        console.log(err);
      }
    });
  }

  allowOnlyNumbers(event: any) {
    const input = event.target;
    input.value = input.value.replace(/[^0-9]/g, '');
    this.sweaterNumber = input.value;
  }

  postPlayer() {
    this.playerService.createPlayer(this.firstName, this.lastName, parseInt(this.sweaterNumber), this.position,
      this.handedness, this.birthCountry, this.dob, this.current_team_id)
      .subscribe({
        next: (responseData) => {
          this.successMessage.set(
            `Player ${responseData.player.first_name} ${responseData.player.last_name} Created`
          );

          setTimeout(() => {
            this.isFading.set(true);
          }, 2500);

          setTimeout(() => {
            this.successMessage.set('');
            this.isFading.set(false);
          }, 2750);
        },
        error: (err) => {
          this.errorMessage.set(
            `Failed to Create Player`
          );

          setTimeout(() => {
            this.isFading.set(true);
          }, 2500);

          setTimeout(() => {
            this.errorMessage.set('');
            this.isFading.set(false);
          }, 2750);
        }
      })
  }

}
