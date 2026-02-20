import { Component, signal, WritableSignal } from '@angular/core';
import { PlayerService } from '../../services/player-service';
import { ActivatedRoute, RouterModule } from '@angular/router';

@Component({
  selector: 'app-view-player',
  imports: [RouterModule],
  templateUrl: './view-player.html',
  styleUrl: './view-player.css',
})
export class ViewPlayer {

  first_name: WritableSignal<string> = signal("");
  last_name: WritableSignal<string> = signal("");
  sweater_number: WritableSignal<number> = signal(-1);
  position: WritableSignal<string> = signal("");
  birth_date: WritableSignal<string> = signal("");
  birth_country: WritableSignal<string> = signal("");
  shoots_catches: WritableSignal<string> = signal("");
  team_full_name: WritableSignal<string> = signal("");
  team_short_name: WritableSignal<string> = signal("");
  avatarUrl: WritableSignal<string> = signal("https://a.espncdn.com/combiner/i?img=/i/headshots/nhl/players/full/5149125.png&w=350&h=254");

  constructor(private playerService: PlayerService, private route: ActivatedRoute) {}

  ngOnInit(){
    const id = this.route.snapshot.paramMap.get('id');
    const playerId = Number(id);

    this.getPlayer(playerId);
  }

  getPlayer(id: number) {
    this.playerService.getPlayer(id)
    .subscribe(
      {
        next: (responseData) => {
          this.first_name.set(responseData.player.first_name);
          this.last_name.set(responseData.player.last_name);
          this.sweater_number.set(responseData.player.sweater_number);
          this.position.set(responseData.player.position);
          this.birth_date.set(responseData.player.birth_date);
          this.birth_country.set(responseData.player.birth_country);
          this.shoots_catches.set(responseData.player.shoots_catches);
          this.team_full_name.set(responseData.player.team_full_name);
          this.team_short_name.set(responseData.player.team_short_name);
        },
        error: (err) => {
          console.log(err);
        }
      }
    )
  }

}
