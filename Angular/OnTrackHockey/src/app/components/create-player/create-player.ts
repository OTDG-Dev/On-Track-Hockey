import { HttpClient } from '@angular/common/http';
import { Component, inject } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { first } from 'rxjs';
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
  sweaterNumber: number = 99;
  position: string = "";
  handedness: string = "";
  birthCountry: string = "";
  dob: string = "";

  allowOnlyNumbers(event: any) {
    const input = event.target;
    input.value = input.value.replace(/[^0-9]/g, '');
    this.sweaterNumber = input.value;
  }
  private http = inject(HttpClient)

  postPlayer() {

    console.log(
      this.firstName,
      this.lastName,
      this.sweaterNumber,
      this.position,
      this.handedness,
      this.birthCountry,
      this.dob
    );

    this.playerService.createPlayer(this.firstName, this.lastName, this.sweaterNumber, this.position, 
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
