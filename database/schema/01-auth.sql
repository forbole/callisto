CREATE TABLE account
(
    address TEXT NOT NULL PRIMARY KEY
);

CREATE TABLE vesting_account
(
    address             TEXT        NOT NULL REFERENCES account (address) PRIMARY KEY,
    original_vesting    JSONB       NOT NULL DEFAULT '{}',
    end_time            TIMESTAMP   NOT NULL,
    start_time          TIMESTAMP   NOT NULL,
    vesting_periods     JSONB       NOT NULL DEFAULT '{}'
);