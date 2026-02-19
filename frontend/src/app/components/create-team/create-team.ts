import { Component, signal, WritableSignal } from '@angular/core';
import { TeamService } from '../../services/team-service';
import { FormsModule } from '@angular/forms';
import { DivisionService } from '../../services/division-service';
import { DivisionData } from '../../interfaces/division-data';

@Component({
  selector: 'app-create-team',
  imports: [FormsModule],
  templateUrl: './create-team.html',
  styleUrl: './create-team.css',
})
export class CreateTeam {

  full_name: string = "";
  short_name: string = "";
  is_active: boolean = false;
  division_id: number = 1;

  divisions: WritableSignal<DivisionData[]> = signal([]);

  successMessage: WritableSignal<string> = signal('');
  errorMessage: WritableSignal<string> = signal('');
  isFading = signal(false);

  constructor(private teamService: TeamService, private divisionService: DivisionService) { }

  ngOnInit()
  {
    this.divisionService.getDivisions()
    .subscribe(
      {
        next: (responseData) => {
          this.divisions.set(responseData.divisions);
        },
        error: (err) => {
          console.log(err);
        }
      }
    )
  }

  postTeam() {
    this.teamService.createTeam(this.full_name, this.short_name, this.is_active, this.division_id)
    .subscribe({
        next: (responseData) => {
            this.successMessage.set(`Team ${responseData.team.full_name} Created`);
  
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
