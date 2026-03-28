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
  
}
