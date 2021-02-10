CREATE TABLE token_price
(
    denom      TEXT      NOT NULL,
    price      NUMERIC   NOT NULL,
    market_cap NUMERIC   NOT NULL,
    timestamp  TIMESTAMP NOT NULL,
    UNIQUE (denom, timestamp)
);
