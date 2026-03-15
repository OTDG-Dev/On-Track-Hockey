import { PlayerData } from "./player-data";

export interface RosterData {
    forwards: PlayerData[];
    defensemen: PlayerData[];
    goalies: PlayerData[];
}
