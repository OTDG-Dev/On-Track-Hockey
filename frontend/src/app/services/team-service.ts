import { Injectable } from '@angular/core';
import { environment } from '../../environments/environment';
import { HttpClient } from '@angular/common/http';
import { TeamData } from '../interfaces/team-data';

@Injectable({
  providedIn: 'root',
})
export class TeamService {

  private baseUrl = environment.apiUrl + "/v1/teams"

  constructor(private httpClient: HttpClient) {}

  getTeams() 
  {
    return this.httpClient.get<{ teams: TeamData[] }>(this.baseUrl);
  }

  createTeam(name: string, short_name: string, is_active: boolean, division_id: string, league_id: string)
  {
    const body = {
      "name": name,
      "short_name": short_name,
      "is_active": is_active,
      "division_id": division_id,
      "league_id": league_id
    };

    return this.httpClient.post<{team: TeamData}>(this.baseUrl, body);
  }
  
}
