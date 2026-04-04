CREATE TABLE IF NOT EXISTS games (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),

    home_team_id BIGINT NOT NULL
        REFERENCES teams(id)
        ON DELETE RESTRICT,

    away_team_id BIGINT NOT NULL
        REFERENCES teams(id)
        ON DELETE RESTRICT,
        
    start_time TIMESTAMP WITH TIME ZONE NOT NULL,

    is_finished BOOLEAN NOT NULL DEFAULT false,

    version integer NOT NULL DEFAULT 1
);