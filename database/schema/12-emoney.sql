/* ---- INFLATION ---- */
CREATE TABLE emoney_inflation
(
    one_row_id              BOOLEAN     NOT NULL DEFAULT TRUE PRIMARY KEY,
    inflation               JSONB       NOT NULL,
    last_applied_time       TIMESTAMP   NOT NULL,
    last_applied_height     BIGINT      NOT NULL,
    height                  BIGINT      NOT NULL,
    CHECK (one_row_id)
);

/* ---- AUTHORITY / GAS PRICE ---- */
CREATE TABLE emoney_gas_prices
(
    authority_key   TEXT        NOT NULL PRIMARY KEY,
    gas_prices      COIN[]      NOT NULL,
    height          BIGINT      NOT NULL
);
CREATE INDEX emoney_gas_prices_index ON emoney_gas_prices (height);