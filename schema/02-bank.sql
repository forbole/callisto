CREATE TABLE supply
(
    coins  COIN[] NOT NULL,
    height BIGINT NOT NULL UNIQUE PRIMARY KEY
);
CREATE INDEX supply_height_index ON supply (height);

CREATE TABLE account_balance
(
    address TEXT   NOT NULL REFERENCES account (address) PRIMARY KEY,
    coins   COIN[] NOT NULL DEFAULT '{}'
);