CREATE TABLE top_accounts
(
    address         TEXT    NOT NULL UNIQUE,
    --  REFERENCES account (address) PRIMARY KEY,
    available       BIGINT,
    delegation      BIGINT,
    redelegation    BIGINT,
    unbonding       BIGINT,
    reward          BIGINT,
    sum             BIGINT  NOT NULL DEFAULT 0
);
CREATE INDEX top_accounts_sum_index ON top_accounts (sum);