CREATE TABLE top_accounts
(
    address         TEXT    NOT NULL REFERENCES account (address) PRIMARY KEY,
    available       BIGINT  NOT NULL DEFAULT 0,
    delegation      BIGINT  NOT NULL DEFAULT 0,
    redelegation    BIGINT  NOT NULL DEFAULT 0,
    unbonding       BIGINT  NOT NULL DEFAULT 0,
    reward          BIGINT  NOT NULL DEFAULT 0,
    sum             BIGINT  NOT NULL DEFAULT 0,
    height          BIGINT  NOT NULL
);
CREATE INDEX top_accounts_height_index ON top_accounts (height);
CREATE INDEX top_accounts_sum_index ON top_accounts (sum);