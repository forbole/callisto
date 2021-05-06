CREATE TABLE token
(
    name TEXT NOT NULL UNIQUE
);

CREATE TABLE token_unit
(
    token_name TEXT NOT NULL REFERENCES token (name),
    denom      TEXT NOT NULL UNIQUE,
    exponent   INT  NOT NULL,
    aliases    TEXT[]
);

CREATE TABLE token_price_history
(
    unit_name  TEXT      NOT NULL REFERENCES token_unit (denom),
    price      NUMERIC   NOT NULL,
    market_cap BIGINT    NOT NULL,
    timestamp  TIMESTAMP NOT NULL,
    CONSTRAINT unique_price_for_timestamp UNIQUE (unit_name, timestamp)
);
CREATE INDEX token_price_timestamp_index ON token_price_history (timestamp);