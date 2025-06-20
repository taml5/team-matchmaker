CREATE TABLE IF NOT EXISTS Player (
    uuid TEXT(128) NOT NULL PRIMARY KEY,
    username TEXT NOT NULL,
    tag TEXT NOT NULL,
    name TEXT NOT NULL,
    UNIQUE (username, tag)
);

CREATE TABLE IF NOT EXISTS Elo (
    player_id TEXT(128) NOT NULL,
    position TEXT NOT NULL CHECK (position IN ('top', 'jgl', 'mid', 'adc', 'sup')),
    elo INTEGER NOT NULL DEFAULT 1000,
    FOREIGN KEY (player_id) REFERENCES Player(uuid) ON DELETE CASCADE
);