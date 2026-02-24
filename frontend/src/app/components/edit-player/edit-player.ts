import { Component, signal, WritableSignal } from '@angular/core';
import { TeamData } from '../../interfaces/team-data';
import { FormsModule } from '@angular/forms';
import { PlayerService } from '../../services/player-service';
import { ActivatedRoute, Router } from '@angular/router';
import { TeamService } from '../../services/team-service';

@Component({
  selector: 'app-edit-player',
  imports: [FormsModule],
  templateUrl: './edit-player.html',
  styleUrl: './edit-player.css',
})
export class EditPlayer {

  playerId: number = -1
  firstName: string = "";
  lastName: string = "";
  sweaterNumber: number = -1
  position: string = "";
  handedness: string = "";
  birthCountry: string = "";
  dob: string = "";
  current_team_id: number = -1;

  teams: WritableSignal<TeamData[]> = signal([]);

  successMessage: WritableSignal<string> = signal('');
  errorMessage: WritableSignal<string> = signal('');
  isFading = signal(false);

  constructor(private playerService: PlayerService, private teamService: TeamService, private route: ActivatedRoute, private router: Router){}

  ngOnInit() {
    const id = this.route.snapshot.paramMap.get('id');
    this.playerId = Number(id);

    this.getPlayer(this.playerId);

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
    this.sweaterNumber = Number(input.value);
  }

  getPlayer(id: number) {
    this.playerService.getPlayer(id)
    .subscribe(
      {
        next: (responseData) => {
          this.firstName = responseData.player.first_name;
          this.lastName = responseData.player.last_name;
          this.sweaterNumber = responseData.player.sweater_number;
          this.position = responseData.player.position;
          this.dob = responseData.player.birth_date;
          this.birthCountry = responseData.player.birth_country;
          this.handedness = responseData.player.shoots_catches;
          this.current_team_id = responseData.player.current_team_id;
        },
        error: (err) => {
          console.log(err);
          this.router.navigate(['/view-players']);
        }
      }
    )
  }

  patchPlayer() {
    this.playerService.patchPlayer(this.firstName, this.lastName, this.sweaterNumber, this.position,
      this.handedness, this.birthCountry, this.dob, this.current_team_id, this.playerId)
      .subscribe({
        next: (responseData) => {
          this.successMessage.set(
            `Player ${responseData.player.first_name} ${responseData.player.last_name} Edited`
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
            `Failed to Edit Player`
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
