import { GameEvent } from "./game-event";

export interface GameData {
    home_team: string,
    away_team: string,
    home_team_id: number,
    away_team_id: number,
    start_time: string,
    game_events: GameEvent[]
}
