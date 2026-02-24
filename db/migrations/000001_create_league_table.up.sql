CREATE TABLE IF NOT EXISTS leagues (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    name text NOT NULL,
    version integer NOT NULL DEFAULT 1
)