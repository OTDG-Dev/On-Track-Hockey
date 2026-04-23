import { SkaterStats } from "./skater-stats";

export interface PlayerData {
  id: number;
  is_active: boolean;
  current_team_id: number;
  first_name: string;
  last_name: string;
  sweater_number: number;
  position: string;
  birth_date: string;
  birth_country: string;
  shoots_catches: string;
  team_full_name: string;
  team_short_name: string;
  skater_stats: SkaterStats;
}
