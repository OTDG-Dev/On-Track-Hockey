import { Component, computed, signal, WritableSignal } from '@angular/core';
import { GameEvent } from '../../interfaces/game-event';
import { GameService } from '../../services/game-service';
import { ActivatedRoute, RouterLink, RouterModule } from '@angular/router';
import { CommonModule, DatePipe } from '@angular/common';
import { FormsModule } from '@angular/forms';
import { PlayerData } from '../../interfaces/player-data';
import { RosterService } from '../../services/roster-service';
import { RosterData } from '../../interfaces/roster-data';

@Component({
  selector: 'app-view-game',
  imports: [CommonModule, DatePipe, FormsModule, RouterLink, RouterModule],
  templateUrl: './view-game.html',
  styleUrl: './view-game.css',
})
export class ViewGame {

  game_id: number = -1

  isAddingEvent = signal(false);

  newEvent = signal({
    event_number: 0,
    period: 1,
    clock_seconds: 0,
    event_type: '',
    situation: '',
    team_id: -1,
    minutes: 0,
    seconds: 0,
    player: ""
  });

  home_team: WritableSignal<string> = signal("");
  away_team: WritableSignal<string> = signal("");
  home_team_id: WritableSignal<number> = signal(-1);
  away_team_id: WritableSignal<number> = signal(-1);
  start_time: WritableSignal<string> = signal("");
  game_events: WritableSignal<GameEvent[]> = signal([]);
  is_finished: WritableSignal<boolean> = signal(false);
  home_roster: WritableSignal<RosterData | null> = signal(null);
  away_roster: WritableSignal<RosterData | null> = signal(null);
  all_players = computed<PlayerData[]>(() => {
    const home = this.home_roster();
    const away = this.away_roster();

    const players = [
      ...(home?.forwards ?? []),
      ...(home?.defensemen ?? []),
      ...(home?.goalies ?? []),
      ...(away?.forwards ?? []),
      ...(away?.defensemen ?? []),
      ...(away?.goalies ?? [])
    ];

    return players.sort((a, b) =>
      a.last_name.localeCompare(b.last_name)
    );
  });


  errorMessage: WritableSignal<string> = signal('');
  completeOrEditMessage: WritableSignal<string> = signal('');
  isFading = signal(false);

  constructor(private gameService: GameService, private rosterService: RosterService, private route: ActivatedRoute) {}

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
        this.is_finished.set(responseData.game.is_finished);

        this.getPlayers();

        console.log(this.game_events);
      },
      error: (err) => {
        console.log(err);
      }
    })
  }

  getPlayers()
  {
    this.rosterService.getRoster(this.home_team_id())
    .subscribe(
      {
        next: (responseData) => 
        {
          this.home_roster.set(responseData.roster);
          console.log(this.home_roster());
        },
        error: (err) => {
          console.log(err);
        }
      }
    )

    this.rosterService.getRoster(this.away_team_id())
    .subscribe(
      {
        next: (responseData) => 
        {
          this.away_roster.set(responseData.roster);
          console.log(this.away_roster());
        },
        error: (err) => {
          console.log(err)
        }
      }
    )

  }

  onAddEvent() {
    this.isAddingEvent.set(true);
  
    this.newEvent.set({
      event_number: this.game_events().length + 1,
      period: 1,
      clock_seconds: 0,
      event_type: '',
      situation: '',
      team_id: this.home_team_id(),
      minutes: 0,
      seconds: 0,
      player: ''
    });
  }
  
  onCancelAdd() {
    this.isAddingEvent.set(false);
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

  onSaveEvent() {
    const minutes = Number(this.newEvent().minutes) || 0;
    const seconds = Number(this.newEvent().seconds) || 0;
  
    const totalSeconds = Number(minutes) * 60 + Number(seconds);
  
    const savedEvent = {
      ...this.newEvent(),
      clock_seconds: totalSeconds,
      team_id: Number(this.newEvent().team_id),
      id: Date.now()
    };

    this.gameService.postGameEvent(Number(savedEvent.period), Number(savedEvent.clock_seconds), savedEvent.event_type, Number(savedEvent.team_id), savedEvent.situation, Number(this.game_id))
    .subscribe(
      {
        next: (responseData) => {    
          this.game_events.set([...this.game_events(), savedEvent]);
          console.log(responseData);
        },
        error: (err) => {
          console.log(err);

          const errorObj = err?.error?.error;
        
          if (typeof errorObj === 'object') {
            const messages = Object.values(errorObj);
            this.errorMessage.set(messages.join(', '));
          } else {
            this.errorMessage.set(errorObj || 'Something went wrong');
          }

          setTimeout(() => {
            this.isFading.set(true);
          }, 2500);

          setTimeout(() => {
            this.errorMessage.set('');
            this.isFading.set(false);
          }, 2750);
        }
      }
    )
    this.isAddingEvent.set(false);
  }

  onMarkComplete() {
    this.gameService.patchGameComplete(Number(this.game_id))
      .subscribe({
        next: () => {
          this.completeOrEditMessage.set('Completed');
          this.is_finished.set(true);
        },
        error: (err) => {
          console.log(err);
        }
      });
  }

  onMarkIncomplete() {
    this.gameService.patchGameIncomplete(Number(this.game_id))
    .subscribe({
      next: () => {
        this.is_finished.set(false);
      },
      error: (err) => {
        console.log(err);
      }
    })
  }

}
