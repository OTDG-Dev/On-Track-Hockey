export interface GameEvent {
    id: number,
    event_number: number,
    period: number,
    clock_seconds: number,
    event_type: string,
    situation: string,
    team_id: number
}
