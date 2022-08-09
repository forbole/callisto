CREATE TABLE marker_params
(
    one_row_id BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    params     JSONB   NOT NULL,
    height     BIGINT  NOT NULL,
    CHECK (one_row_id)
);

CREATE TABLE marker_account
(
    address                  TEXT       NOT NULL PRIMARY KEY REFERENCES account (address),
    access_control           TEXT       NOT NULL,
    allow_governance_control BOOLEAN    NOT NULL,
    denom                    TEXT       NOT NULL,
    marker_type              TEXT       NOT NULL,
    status                   TEXT       NOT NULL,
    total_supply             TEXT       NOT NULL,
    height                   BIGINT     NOT NULL
);
CREATE INDEX marker_account_height_index ON marker_account (height);
CREATE INDEX marker_account_address_index ON marker_account (address);

CREATE TABLE marker_token_price
(
    id         SERIAL                      NOT NULL PRIMARY KEY,
    unit_name  TEXT                        NOT NULL UNIQUE REFERENCES marker_account (denom),
    price      DECIMAL                     NOT NULL,
    market_cap BIGINT                      NOT NULL,
    timestamp  TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
CREATE INDEX marker_token_price_timestamp_index ON marker_token_price (timestamp);
