import { Injectable } from '@angular/core';
import { environment } from '../../environments/environment';
import { HttpClient } from '@angular/common/http';
import { DivisionData } from '../interfaces/division-data';

@Injectable({
  providedIn: 'root',
})
export class DivisionService {

  private baseUrl = environment.apiUrl + "/v1/divisions";

  constructor(private httpClient: HttpClient) {}

  getDivisions() {
    return this.httpClient.get<{ divisions: DivisionData[] }>(this.baseUrl);
  }

  createDivision(league_id: number, name: string)
  {
    const body = {
      "league_id": league_id,
      "name": name,
    };

    return this.httpClient.post<{ division: DivisionData }>(this.baseUrl, body);
  }
  
}
