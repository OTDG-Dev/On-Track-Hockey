import { Injectable } from '@angular/core';
import { environment } from '../../environments/environment';
import { HttpClient } from '@angular/common/http';
import { GameData } from '../interfaces/game-data';

@Injectable({
  providedIn: 'root',
})
export class GameService {

  private baseUrl = environment.apiUrl + "/v1/games"

  constructor(private httpClient: HttpClient) { }

  getGame(id: number) {
    return this.httpClient.get<{ game: GameData }>(`${this.baseUrl}/${id}`);
  }

  getGames(){
    return this.httpClient.get<{ games: GameData[]}>(`${this.baseUrl}`)
  }

  postGameEvent(period: number, clock_seconds: number, event_type: string, team_id: number, situation: string, game_id: number) {
    
    const body = {
      "period": period,
      "clock_seconds": clock_seconds,
      "event_type": event_type,
      "team_id": team_id,
      "situation": situation
    }

    return this.httpClient.post<{ gameEvents: GamepadEvent[]}>(`${this.baseUrl}/${game_id}/events`, body);
  }
  
}
