CREATE TABLE IF NOT EXISTS event_participants (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,

    game_event_id BIGINT NOT NULL
        REFERENCES game_events(id)
        ON DELETE CASCADE,

    player_id BIGINT NOT NULL
        REFERENCES players(id)
        ON DELETE RESTRICT,

    role TEXT NOT NULL,

    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_event_participants_game_event_id ON event_participants(game_event_id);
CREATE INDEX idx_event_participants_player_id ON event_participants(player_id);
