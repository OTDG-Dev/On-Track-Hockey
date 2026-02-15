CREATE TABLE IF NOT EXISTS players (
    id bigint PRIMARY KEY GENERATED ALWAYS AS IDENTITY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    is_active boolean NOT NULL DEFAULT false,
    current_team_id bigint NOT NULL,
    first_name text NOT NULL,
    last_name text NOT NULL,
    sweater_number smallint NOT NULL
        CHECK (sweater_number BETWEEN 0 AND 99),
    position text NOT NULL
        CHECK (position IN ('C', 'LW', 'RW', 'D', 'G')),
    birth_date DATE NOT NULL,
    birth_country text NOT NULL,
    headshot text NOT NULL,
    shoots_catches text NOT NULL
        CHECK (shoots_catches IN ('L', 'R')),
    version integer NOT NULL DEFAULT 1
)