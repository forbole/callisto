CREATE TABLE supply
(
    coins  COIN[] NOT NULL,
    height BIGINT NOT NULL UNIQUE PRIMARY KEY
);
CREATE INDEX supply_height_index ON supply (height);

CREATE TABLE account_balance
(
    address TEXT   NOT NULL REFERENCES account (address),
    coins   COIN[] NOT NULL DEFAULT '{}',
    height  BIGINT NOT NULL,
    PRIMARY KEY (address, height)
);
CREATE INDEX account_balance_address_index ON account_balance (address);
CREATE INDEX account_balance_height_index ON account_balance (height);