CREATE TABLE players (
    id bigserial primary key,
    name varchar NOT NULL,

    UNIQUE(name)
);

CREATE TABLE players_scores (
    id bigserial primary key,
    score bigint,
    player_id bigint NOT NULL REFERENCES players,
    updated_at timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,

    CONSTRAINT unique_player_fk UNIQUE(player_id)
);
