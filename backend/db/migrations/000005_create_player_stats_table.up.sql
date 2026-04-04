CREATE TABLE IF NOT EXISTS player_stats (
    player_id bigint PRIMARY KEY
        REFERENCES players(id)
        ON DELETE CASCADE,
    skater_stats jsonb,
    goalie_stats jsonb,
    CHECK (num_nonnulls(skater_stats, goalie_stats) = 1)
);