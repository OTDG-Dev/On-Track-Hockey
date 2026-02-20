import { Component } from '@angular/core';
import { PlayerService } from '../../services/player-service';
import { ActivatedRoute, RouterModule } from '@angular/router';

@Component({
  selector: 'app-view-player',
  imports: [RouterModule],
  templateUrl: './view-player.html',
  styleUrl: './view-player.css',
})
export class ViewPlayer {

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
          console.log(responseData.player);
        },
        error: (err) => {
          console.log(err);
        }
      }
    )
  }

}
