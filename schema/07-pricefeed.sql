CREATE TABLE token
(
    name        TEXT NOT NULL UNIQUE,
    traded_unit TEXT NOT NULL UNIQUE
);

CREATE TABLE token_unit
(
    token_name TEXT   NOT NULL REFERENCES token (name),
    denom      TEXT   NOT NULL UNIQUE,
    exponent   INT    NOT NULL,
    aliases    TEXT[] NOT NULL DEFAULT ARRAY []::TEXT[]
);

CREATE TABLE token_price
(
    /* Needed for the below token_price function to work properly */
    id         SERIAL    NOT NULL PRIMARY KEY,

    name       TEXT      NOT NULL REFERENCES token_unit (denom),
    price      NUMERIC   NOT NULL,
    market_cap BIGINT    NOT NULL,
    timestamp  TIMESTAMP NOT NULL,
    UNIQUE (name, timestamp)
);

/**
  * This function is used to have a Hasura compute field (https://hasura.io/docs/1.0/graphql/core/schema/computed-fields.html)
  * inside the account_balance table, so that it's easy to determine the token price that is associated with that balance.
 */
CREATE FUNCTION token_price(account_balance_row account_balance) RETURNS SETOF token_price AS
$$
SELECT id, name, price, market_cap, timestamp
FROM (
         SELECT DISTINCT ON (name) name, id, price, market_cap, timestamp
         FROM (
                  SELECT *
                  FROM token_price
                  WHERE timestamp <= (SELECT timestamp FROM block WHERE block.height = account_balance_row.height)
                  ORDER BY timestamp DESC
              ) AS prices
     ) as prices
$$ LANGUAGE sql STABLE;