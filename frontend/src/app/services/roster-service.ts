import { Injectable } from '@angular/core';
import { environment } from '../../environments/environment';
import { HttpClient } from '@angular/common/http';
import { RosterData } from '../interfaces/roster-data';

@Injectable({
  providedIn: 'root',
})
export class RosterService {

  private baseUrl = environment.apiUrl + "/v1/roster"

  constructor(private httpClient: HttpClient) { }

  getRoster(id: number)
  {
    return this.httpClient.get<{ roster: RosterData}>(`${this.baseUrl}/${id}`);
  }
  
}
