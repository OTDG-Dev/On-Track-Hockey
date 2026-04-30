import { GameEvent } from "./game-event";

export interface GameData {
    id: number,
    home_team: string,
    away_team: string,
    home_team_id: number,
    away_team_id: number,
    start_time: string,
    is_finished: boolean,
    game_events: GameEvent[]
}
