import { Component, signal, WritableSignal } from '@angular/core';
import { GameEvent } from '../../interfaces/game-event';
import { GameService } from '../../services/game-service';
import { ActivatedRoute } from '@angular/router';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-view-game',
  imports: [CommonModule],
  templateUrl: './view-game.html',
  styleUrl: './view-game.css',
})
export class ViewGame {

  game_id: number = -1

  home_team: WritableSignal<string> = signal("");
  away_team: WritableSignal<string> = signal("");
  home_team_id: WritableSignal<number> = signal(-1);
  away_team_id: WritableSignal<number> = signal(-1);
  start_time: WritableSignal<string> = signal("");
  game_events: WritableSignal<GameEvent[]> = signal([]);

  constructor(private gameService: GameService, private route: ActivatedRoute) {}

  ngOnInit()
  {
    const id = this.route.snapshot.paramMap.get('id');
    this.game_id = Number(id);

    this.getGame(this.game_id);
  }

  getGame(id: number)
  {
    this.gameService.getGame(id)
    .subscribe({
      next: (responseData) => {
        this.home_team.set(responseData.game.home_team);
        this.away_team.set(responseData.game.away_team);
        this.home_team_id.set(responseData.game.home_team_id);
        this.away_team_id.set(responseData.game.away_team_id);
        this.start_time.set(responseData.game.start_time);
        this.game_events.set(responseData.game.game_events);
        console.log(this.game_events);
      },
      error: (err) => {
        console.log(err);
      }
    })
  }

  trackByEvent(index: number, event: any) {
    return event.id;
  }
  
  getTeamName(teamId: number): string {
    if (teamId === this.home_team_id()) return this.home_team();
    if (teamId === this.away_team_id()) return this.away_team();
    return 'Unknown';
  }
  
  formatClock(seconds: number): string {
    const mins = Math.floor(seconds / 60);
    const secs = seconds % 60;
    return `${mins}:${secs.toString().padStart(2, '0')}`;
  }

}
