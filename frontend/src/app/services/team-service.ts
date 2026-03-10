import { Injectable } from '@angular/core';
import { environment } from '../../environments/environment';
import { HttpClient } from '@angular/common/http';
import { TeamData } from '../interfaces/team-data';

@Injectable({
  providedIn: 'root',
})
export class TeamService {

  private baseUrl = environment.apiUrl + "/v1/teams"

  constructor(private httpClient: HttpClient) { }

  getTeam(id: number) {
    return this.httpClient.get<{ team: TeamData}>(`${this.baseUrl}/${id}`);
  }

  getTeams() {
    return this.httpClient.get<{ teams: TeamData[] }>(this.baseUrl);
  }

  createTeam(full_name: string, short_name: string, is_active: boolean, division_id: number) {
    const body = {
      "full_name": full_name,
      "short_name": short_name,
      "is_active": is_active,
      "division_id": division_id
    };

    return this.httpClient.post<{ team: TeamData }>(this.baseUrl, body);
  }

  patchTeam(full_name: string, short_name: string, is_active: boolean, division_id: number, id: number) {
    const body = {
      "full_name": full_name,
      "short_name": short_name,
      "is_active": is_active,
      "division_id": division_id
    };

    return this.httpClient.patch<{ team: TeamData }>(`${this.baseUrl}/${id}`, body);
  }

  deactivateTeam(id: number) {
    const body = {
      "is_active": false
    };
    return this.httpClient.patch<{ team: TeamData}>(`${this.baseUrl}/${id}`, body);
  }

  deleteTeam(id: number) {
    return this.httpClient.delete<{ team: TeamData}>(`${this.baseUrl}/${id}`);
  }

}
