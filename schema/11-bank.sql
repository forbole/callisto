CREATE TABLE account_balance
(
    address TEXT   NOT NULL REFERENCES account (address),
    coins   COIN[] NOT NULL DEFAULT '{}',
    height  BIGINT NOT NULL,
    PRIMARY KEY (address, height)
);