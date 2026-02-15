import { Component, signal, WritableSignal } from '@angular/core';
import { PlayerService } from '../../services/player-service';
import { PlayerData } from '../../interfaces/player-data';
import { CommonModule } from '@angular/common';

@Component({
  selector: 'app-view-players',
  standalone: true,
  imports: [CommonModule],
  templateUrl: './view-players.html',
  styleUrl: './view-players.css',
})
export class ViewPlayers {

  players: WritableSignal<PlayerData[]> = signal([]);

  constructor(private playerService: PlayerService) { }

  ngOnInit() {
    console.log('Before request:', this.players);

    this.playerService.getPlayers().subscribe({
      next: (responseData) => {
        this.players.set(responseData.players);
      },
      error: (err) => {
        console.log(err);
      }
    });
  }

  onPositionChange(position: string) {
    this.playerService.getPlayers(position).subscribe({
      next: (responseData) => {
        this.players.set(responseData.players);
      },
      error: (err) => {
        console.error(err);
      }
    });
  }

}
