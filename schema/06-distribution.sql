CREATE TYPE DEC_COIN AS
(
    denom  TEXT,
    amount NUMERIC
);

CREATE TABLE community_pool
(
    coins                  DEC_COIN[]               NOT NULL,
    height                BIGINT                NOT NULL,
    PRIMARY KEY (coins,height)
);