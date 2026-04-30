import { Component, computed, signal, WritableSignal } from '@angular/core';
import { EventService } from '../../services/event-service';
import { ActivatedRoute, RouterLink, RouterModule } from '@angular/router';
import { TitleCasePipe } from '@angular/common';
import { TeamService } from '../../services/team-service';
import { FormsModule } from '@angular/forms';
import { PlayerData } from '../../interfaces/player-data';
import { ParticipantData } from '../../interfaces/participant-data';
import { RosterService } from '../../services/roster-service';
import { RosterData } from '../../interfaces/roster-data';
import { PlayerService } from '../../services/player-service';

@Component({
  selector: 'app-view-event',
  imports: [FormsModule ,RouterModule, TitleCasePipe],
  templateUrl: './view-event.html',
  styleUrl: './view-event.css',
})
export class ViewEvent {

  event_id: number = -1;

  id: WritableSignal<number> = signal(-1);
  event_number: WritableSignal<number> = signal(-1);
  period: WritableSignal<number> = signal(-1);
  clock_seconds: WritableSignal<number> = signal(-1);
  event_type: WritableSignal<string> = signal("");
  situation: WritableSignal<string> = signal("");
  team_id: number = -1;
  participants: WritableSignal<ParticipantData[]> = signal([]);
  team_name:  WritableSignal<string> = signal("");
  roster: WritableSignal<RosterData | null> = signal(null);
  all_players = computed<PlayerData[]>(() => {
    const home = this.roster();

    const players = [
      ...(home?.forwards ?? []),
      ...(home?.defensemen ?? []),
      ...(home?.goalies ?? []),
    ];

    return players.sort((a, b) =>
      a.last_name.localeCompare(b.last_name)
    );
  });

  newParticipant = signal({
    id: -1,
    role: "",
    event_id: -1,
    player_id: -1
  });

  isAddingParticipant = signal(false);

  constructor(private eventService: EventService, private teamService: TeamService, private rosterService: RosterService, private playerService: PlayerService, private route: ActivatedRoute){}

  ngOnInit()
  {
    const id = this.route.snapshot.paramMap.get('id');
    this.event_id = Number(id);

    this.getEvent(this.event_id);
  }

  getEvent(id: number) {
    this.eventService.getEvent(id)
    .subscribe(
      {
        next: (responseData) => 
        {
          this.id.set(responseData.game_events.id);
          this.event_number.set(responseData.game_events.event_number);
          this.period.set(responseData.game_events.period);
          this.clock_seconds.set(responseData.game_events.clock_seconds);
          this.event_type.set(responseData.game_events.event_type);
          this.situation.set(responseData.game_events.situation);
          this.team_id = responseData.game_events.team_id;

          this.getParticipants(id)
          this.teamService.getTeam(this.team_id)
          .subscribe(
            {
              next: (teamResponse) => {
                this.team_name.set(teamResponse.team.full_name);
              },
              error: (err) => {
                console.log(err);
              }
            }
          );
          this.getRoster(this.team_id);

        },
        error: (err) => {
          console.log(err);
        }
      }
    )
  }

  getParticipants(id: number){
    this.eventService.getParticipants(id)
    .subscribe(
      {
        next: (responseData) =>
        {
          this.participants.set(responseData.game_event_participants ?? []);
          console.log(responseData);
          console.log(this.participants());
        },
        error: (err) => 
        {
          console.log(err);
        }
      }
    )
  }

  getRoster(team_id: number)
  {
    this.rosterService.getRoster(team_id)
    .subscribe(
      {
        next: (responseData) =>
        {
          this.roster.set(responseData.roster);
        },
        error: (err) =>
        {
          console.log(err);
        }
      }
    )
  }

  getPlayer(player_id: number)
  {
    const player = this.all_players().find(p => p.id === player_id);
    return player ? `${player.first_name} ${player.last_name}` : 'Unknown';
  }

  onAddParticipant(){
    this.isAddingParticipant.set(true);
  }

  clockDisplay = computed(() => {
    const total = this.clock_seconds();
    if (total < 0) return '';
  
    const minutes = Math.floor(total / 60);
    const seconds = total % 60;
  
    return `${minutes}:${seconds.toString().padStart(2, '0')}`;
  });

  onCancelParticipant() {
    this.isAddingParticipant.set(false);
  }

  onSaveParticipant() {
    const savedParticipant = {
      ...this.newParticipant()
    }

    this.eventService.postParticipant(this.event_id, this.newParticipant().role, this.newParticipant().player_id)
    .subscribe(
      {
        next: (responseData) => {
          this.participants.set([...this.participants(), savedParticipant]);
          console.log(responseData);
        },
        error: (err) =>
        {
          console.log(err);
        }
      }
    )
    this.isAddingParticipant.set(false);
  }

  addParticipant() {
    throw new Error('Method not implemented.');
  }
  removeParticipant(_t41: any) {
    throw new Error('Method not implemented.');
  }

}
