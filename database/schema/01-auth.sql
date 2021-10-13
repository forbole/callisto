CREATE TABLE account
(
    address TEXT NOT NULL PRIMARY KEY
);

CREATE TYPE VESTING_PERIOD
(
    length      TEXT    NOT NULL,
    amounts     COIN[]  NOT NULL DEFAULT '{}'
);

CREATE TABLE vesting_account
(
    address             TEXT                NOT NULL REFERENCES account (address) PRIMARY KEY,
    original_vesting    COIN[]              NOT NULL DEFAULT '{}',
    end_time            TIMESTAMP           NOT NULL,
    start_time          TIMESTAMP           NOT NULL,
    vesting_periods     VESTING_PERIOD[]    NOT NULL DEFAULT '{}'
);