CREATE TYPE COIN AS
(
    denom  TEXT,
    amount TEXT
);

CREATE TABLE account_balance_history
(
    address TEXT   NOT NULL REFERENCES account (address),
    coins   COIN[] NOT NULL DEFAULT '{}',
    height  BIGINT NOT NULL,
    CONSTRAINT unique_balance_for_height UNIQUE (address, height)
);
CREATE INDEX account_balance_height_index ON account_balance_history (height);

