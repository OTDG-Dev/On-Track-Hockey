CREATE TABLE IF NOT EXISTS teams (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    full_name text NOT NULL,
    short_name text NOT NULL
        CHECK (char_length(short_name) = 3),
    division_id bigint NOT NULL -- reference div table later
)