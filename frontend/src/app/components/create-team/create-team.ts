import { Component, signal, WritableSignal } from '@angular/core';
import { TeamService } from '../../services/team-service';
import { FormsModule } from '@angular/forms';

@Component({
  selector: 'app-create-team',
  imports: [FormsModule],
  templateUrl: './create-team.html',
  styleUrl: './create-team.css',
})
export class CreateTeam {

  name: string = "";
  short_name: string = "";
  is_active: boolean = false;
  division_id: string = "";
  league_id: string = "";

  successMessage: WritableSignal<string> = signal('');
  errorMessage: WritableSignal<string> = signal('');
  isFading = signal(false);

  constructor(private teamService: TeamService) {}

  postTeam() {
    this.teamService.createTeam(this.name, this.short_name, this.is_active, this.division_id, this.league_id)
    .subscribe({
      next: (responseData) => {
        this.successMessage.set(
          `Team ${responseData.team.name} Created`
        );

        setTimeout(() => {
          this.isFading.set(true);
        }, 2500);
      
        setTimeout(() => {
          this.successMessage.set('');
          this.isFading.set(false);
        }, 2750);
      },
      error: (err) => {
        this.errorMessage.set(
          `Failed to Create Team`
        );

        setTimeout(() => {
          this.isFading.set(true);
        }, 2500);
      
        setTimeout(() => {
          this.errorMessage.set('');
          this.isFading.set(false);
        }, 2750);
      }
    })
  }
  
}
