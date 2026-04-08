import { Component, signal, WritableSignal } from '@angular/core';
import { GameService } from '../../services/game-service';
import { GameData } from '../../interfaces/game-data';
import { DatePipe } from '@angular/common';
import { RouterLink } from "@angular/router";

@Component({
  selector: 'app-view-games',
  imports: [DatePipe, RouterLink],
  templateUrl: './view-games.html',
  styleUrl: './view-games.css',
})
export class ViewGames {

  games: WritableSignal<GameData[]> = signal([]);

  constructor(private gameService: GameService) {}

  ngOnInit(){

    this.gameService.getGames()
    .subscribe(
      {
        next: (responseData) => {
          this.games.set(responseData.games);
        },
        error: (err) => {
          console.log(err);
        }
      }
    )
  }

}
