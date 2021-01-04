CREATE TYPE COIN AS
(
    denom  TEXT,
    amount BIGINT
);

CREATE TABLE account
(
    address TEXT NOT NULL PRIMARY KEY
);