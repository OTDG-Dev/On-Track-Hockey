import { Injectable } from '@angular/core';
import { environment } from '../../environments/environment';
import { HttpClient } from '@angular/common/http';
import { LeagueData } from '../interfaces/league-data';
import { DivisionData } from '../interfaces/division-data';

@Injectable({
  providedIn: 'root',
})
export class LeagueService {

  private baseUrl = environment.apiUrl + "/v1/leagues";

  constructor(private httpClient: HttpClient) {}

  createLeague(name: string)
  {
    const body = {
      "name": name
    };

    return this.httpClient.post<{ league: LeagueData }>(this.baseUrl, body);
  }

  getLeagues() {
    return this.httpClient.get<{ leagues: LeagueData[] }>(this.baseUrl);
  }

  getDivisionByLeague(id: number) {
    return this.httpClient.get<{ divisions: DivisionData[] }>(`${this.baseUrl}/${id}/divisions`)
  }
  
}
