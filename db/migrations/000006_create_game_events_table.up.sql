CREATE TABLE IF NOT EXISTS game_events (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1,

    game_id BIGINT NOT NULL
        REFERENCES games(id)
        ON DELETE CASCADE,

    period INTEGER NOT NULL,

    clock INTERVAL NOT NULL,

    event_type text NOT NULL
        CHECK (event_type in ('penalty', 'goal', 'shot', 'save')),

    team_id bigint NOT NULL
        REFERENCES teams(id)
        ON DELETE RESTRICT
);