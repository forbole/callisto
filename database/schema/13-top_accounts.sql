CREATE TABLE top_accounts
(
    address         TEXT   NOT NULL REFERENCES account (address) PRIMARY KEY,
    available       BIGINT DEFAULT 0,
    delegation      BIGINT DEFAULT 0,
    redelegation    BIGINT DEFAULT 0,
    unbonding       BIGINT DEFAULT 0,
    reward          BIGINT DEFAULT 0,
    sum             BIGINT NOT NULL DEFAULT 0,
    height          BIGINT NOT NULL
);
CREATE INDEX top_accounts_sum_index ON top_accounts (sum);
CREATE INDEX top_accounts_height_index ON top_accounts (height);

CREATE TABLE top_accounts_params
(
    one_row_id      BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    total_accounts  BIGINT  NOT NULL,
    height          BIGINT  NOT NULL,
    CHECK (one_row_id)
);
CREATE INDEX top_accounts_params_height_index ON top_accounts_params (height);