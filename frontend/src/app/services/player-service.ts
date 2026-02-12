import { HttpClient } from "@angular/common/http";
import { Injectable } from '@angular/core';
import { PlayerData } from "../interfaces/player-data";

@Injectable({
  providedIn: 'root',
})
export class PlayerService {

  constructor(private httpClient: HttpClient) {

  }

  getPlayers() {
    return this.httpClient.get<PlayerData[]>("http://localhost:3000/v1/players");
  }

  createPlayer(firstName: string, lastName: string, sweaterNumber: number, position: string, handedness: string,
    birth_country: string, dob: string) {
    const body = {
      "first_name": firstName,
      "last_name": lastName,
      "sweater_number": sweaterNumber,
      "position": position,
      "birth_date": dob,
      "birth_country": birth_country,
      "shoots_catches": handedness
    }

    return this.httpClient.post<PlayerData>(`http://localhost:3000/v1/players`, body)
  }
}
