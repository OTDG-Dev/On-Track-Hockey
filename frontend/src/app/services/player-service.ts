import { HttpClient } from "@angular/common/http";
import { Injectable } from '@angular/core';
import { PlayerData } from "../interfaces/player-data";
import { environment } from "../../environments/environment";

@Injectable({
  providedIn: 'root',
})
export class PlayerService {

  private baseUrl = environment.apiUrl + "/v1/players";

  constructor(private httpClient: HttpClient) {}

  getPlayers(position?: string) 
  {
    const params: any = {};

    if (position) {
      params.position = position;
    }

    return this.httpClient.get<{ players: PlayerData[] }>(
      this.baseUrl,
      { params }
    );
  }

  getPlayer(id: number)
  {
    return this.httpClient.get<{ player: PlayerData}>(`${this.baseUrl}/${id}`);
  }

  createPlayer(firstName: string, lastName: string, sweaterNumber: number, position: string,
    handedness: string, birth_country: string, dob: string, current_team_id: number) {

    const body = {
      "first_name": firstName,
      "last_name": lastName,
      "sweater_number": sweaterNumber,
      "position": position,
      "birth_date": dob,
      "birth_country": birth_country,
      "shoots_catches": handedness,
      "current_team_id": current_team_id,
    };

    return this.httpClient.post<{ player: PlayerData }>(this.baseUrl, body);
  }

  patchPlayer(firstName: string, lastName: string, sweaterNumber: number, position: string,
    handedness: string, birth_country: string, dob: string, current_team_id: number, id: number){
      
      const body = {
        "first_name": firstName,
        "last_name": lastName,
        "sweater_number": sweaterNumber,
        "position": position,
        "birth_date": dob,
        "birth_country": birth_country,
        "shoots_catches": handedness,
        "current_team_id": current_team_id,
      };

      return this.httpClient.patch<{ player: PlayerData }>(`${this.baseUrl}/${id}`, body);
    }

  deletePlayer(id: number) {
    return this.httpClient.delete<{ player: PlayerData}>(`${this.baseUrl}/${id}`);
  }
}
