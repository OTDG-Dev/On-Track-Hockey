import { Component, signal, WritableSignal } from '@angular/core';
import { DivisionData } from '../../interfaces/division-data';
import { TeamService } from '../../services/team-service';
import { DivisionService } from '../../services/division-service';
import { FormsModule } from '@angular/forms';
import { ActivatedRoute, Router } from '@angular/router';

@Component({
  selector: 'app-edit-team',
  imports: [FormsModule],
  templateUrl: './edit-team.html',
  styleUrl: './edit-team.css',
})
export class EditTeam {

  teamId: number = -1;
  full_name: WritableSignal<string> = signal("");
  short_name: WritableSignal<string> = signal("");
  is_active: WritableSignal<boolean> = signal(false);
  division_id: WritableSignal<number> = signal(-1);

  divisions: WritableSignal<DivisionData[]> = signal([]);

  successMessage: WritableSignal<string> = signal('');
  errorMessage: WritableSignal<string> = signal('');
  isFading = signal(false);

  constructor(private teamService: TeamService, private divisionService: DivisionService, private route: ActivatedRoute, private router: Router) { }

  ngOnInit()
  {
    const id = this.route.snapshot.paramMap.get('id');
    this.teamId = Number(id);

    this.getTeam(this.teamId);
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

  getTeam(id: number)
  {
    this.teamService.getTeam(id)
    .subscribe(
      {
        next: (responseData) => {
          this.full_name.set(responseData.team.full_name);
          this.short_name.set(responseData.team.short_name);
          this.is_active.set(responseData.team.is_active);
          this.division_id.set(responseData.team.division_id);
          console.log("is_active : " + this.is_active())
        },
        error: (err) => {
          console.log(err);
          this.router.navigate(['/view-teams']);
        }
      }
    )
  }

  patchTeam() {
    this.teamService.patchTeam(this.full_name(), this.short_name(), this.is_active(), this.division_id(), this.teamId)
    .subscribe({
        next: (responseData) => {
            this.successMessage.set(`Team ${responseData.team.full_name} Edited`);
  
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
              `Failed to Edit Team`
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
