CREATE TYPE COIN AS
(
    denom  TEXT,
    amount TEXT
);

CREATE TABLE account
(
    address TEXT NOT NULL PRIMARY KEY
);