CREATE TABLE IF NOT EXISTS game_event_participants (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1,

    event_id BIGINT NOT NULL
        REFERENCES game_events(id)
        ON DELETE CASCADE,

    player_id BIGINT NOT NULL
        REFERENCES players(id)
        ON DELETE RESTRICT,

    role text NOT NULL
        CHECK (role IN ('scorer', 'assist_primary', 'assist_secondary','penalty_taker')),
    
    UNIQUE (event_id, role),
    UNIQUE (event_id, player_id)
)