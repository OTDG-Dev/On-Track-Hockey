CREATE TABLE IF NOT EXISTS divisions (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    league_id bigint NOT NULL
        REFERENCES leagues(id)
        ON DELETE RESTRICT,
    version integer NOT NULL DEFAULT 1
)