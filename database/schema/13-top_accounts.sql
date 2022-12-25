CREATE TABLE top_accounts
(
    address         TEXT   NOT NULL REFERENCES account (address) PRIMARY KEY,
    available       BIGINT DEFAULT 0,
    delegation      BIGINT DEFAULT 0,
    redelegation    BIGINT DEFAULT 0,
    unbonding       BIGINT DEFAULT 0,
    reward          BIGINT DEFAULT 0,
    sum             BIGINT NOT NULL DEFAULT 0
);
CREATE INDEX top_accounts_sum_index ON top_accounts (sum);