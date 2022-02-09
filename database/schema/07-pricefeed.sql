/* ---- TOKENS ---- */

CREATE TABLE token
(
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE token_unit
(
    token_name TEXT NOT NULL REFERENCES token (name),
    denom      TEXT NOT NULL UNIQUE,
    exponent   INT  NOT NULL,
    aliases    TEXT[],
    price_id   TEXT
);


/* ---- TOKEN PRICES ---- */

CREATE TABLE token_price
(
    /* Needed for the below token_price function to work properly */
    id         SERIAL                      NOT NULL PRIMARY KEY,

    unit_name  TEXT                        NOT NULL REFERENCES token_unit (denom) UNIQUE,
    price      DECIMAL                     NOT NULL,
    market_cap BIGINT                      NOT NULL,
    timestamp  TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
CREATE INDEX token_price_timestamp_index ON token_price (timestamp);


CREATE TABLE token_price_history
(
    id         SERIAL                      NOT NULL PRIMARY KEY,
    unit_name  TEXT                        NOT NULL REFERENCES token_unit (denom),
    price      DECIMAL                     NOT NULL,
    market_cap BIGINT                      NOT NULL,
    timestamp  TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    CONSTRAINT unique_price_for_timestamp UNIQUE (unit_name, timestamp)
);
CREATE INDEX token_price_history_timestamp_index ON token_price_history (timestamp);
