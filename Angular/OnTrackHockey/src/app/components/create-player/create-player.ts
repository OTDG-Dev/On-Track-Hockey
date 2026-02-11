import { HttpClient } from '@angular/common/http';
import { Component, inject } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { first } from 'rxjs';

@Component({
  selector: 'app-create-player',
  imports: [FormsModule],
  templateUrl: './create-player.html',
  styleUrl: './create-player.css',
})
export class CreatePlayer {

  private http = inject(HttpClient)


  firstName: string = "";
  lastName: string = "";
  sweaterNumber: number | null = null;
  position: string = "";
  handedness: string = "";
  dob: string = "";

  onClick() {
    console.log(
      this.firstName,
      this.lastName,
      this.sweaterNumber,
      this.position,
      this.handedness,
      this.dob
    );

    const payload = {
      "first_name": this.firstName,
      "last_name": this.lastName,
      "sweater_number": this.sweaterNumber,
      "position": this.position,
      "birth_date": "1997-01-13",
      "birth_country": "CAN",
      "shoots_catches": "L"
    }
    console.log(this.http.post<string>('localhost:3000/v1/players', payload))
  }

}
