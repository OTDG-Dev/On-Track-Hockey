import { Component, inject } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { PlayerService } from '../../services/player-service';

@Component({
  selector: 'app-create-player',
  imports: [FormsModule],
  templateUrl: './create-player.html',
  styleUrl: './create-player.css',
})
export class CreatePlayer {

  constructor(private playerService: PlayerService) {}

  firstName: string = "";
  lastName: string = "";
  sweaterNumber: string = "";
  position: string = "";
  handedness: string = "";
  birthCountry: string = "";
  dob: string = "";

  allowOnlyNumbers(event: any) {
    const input = event.target;
    input.value = input.value.replace(/[^0-9]/g, '');
    this.sweaterNumber = input.value;
  }

  postPlayer() {
    this.playerService.createPlayer(this.firstName, this.lastName, parseInt(this.sweaterNumber), this.position, 
                                    this.handedness, this.birthCountry, this.dob)
    .subscribe({
      next: (responseData) => {
        console.log(responseData);
      },
      error: (err) => {
        console.log(err);
      }
    })

  }

}
