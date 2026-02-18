CREATE TABLE IF NOT EXISTS divisions (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    name text NOT NULL,
    league_id bigint NOT NULL
        REFERENCES leagues(id)
        ON DELETE RESTRICT
)