CREATE TABLE account
(
    address TEXT NOT NULL PRIMARY KEY
);

CREATE TYPE PERIOD AS
(
    length  BIGINT,
    amount COIN[]
);

# vestingAccountYAML
CREATE TABLE vestingAccount
(
    address TEXT NOT NULL REFERENCES account (address) PRIMARY KEY,
    pub_key TEXT NOT NULL,
    account_number TEXT NOT NULL,
    sequence TEXT NOT NULL,
    original_vesting COIN[],
    delegated_free   COIN[],
    delegated_vesting COIN[],
    end_time BIGINT,
    StartTime BIGINT,
    VestingPeriods PERIOD[]
);

