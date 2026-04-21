import { Injectable } from '@angular/core';
import { environment } from '../../environments/environment';
import { HttpClient } from '@angular/common/http';
import { GameEvent } from '../interfaces/game-event';
import { ParticipantData } from '../interfaces/participant-data';

@Injectable({
  providedIn: 'root',
})
export class EventService {

  private baseUrl = environment.apiUrl + "/v1/events";

  constructor(private httpClient: HttpClient) {}

  getEvent(id: number){
    return this.httpClient.get<{ game_events: GameEvent  }>(`${this.baseUrl}/${id}`);
  }
  
  getParticipants(id: number){
    return this.httpClient.get<{ game_event_participant: ParticipantData[]  }>(`${this.baseUrl}/${id}/participants`);
  }

  postParticipant(id: number, role: string, player_id: number){
    const body = {
      role: role,
      player_id: player_id
    }

    return this.httpClient.post<{ game_event_participant: ParticipantData[]  }>(`${this.baseUrl}/${id}/participants`, body);
  }
}
