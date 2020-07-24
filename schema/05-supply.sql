CREATE TABLE supply
(
    coins  COIN[] NOT NULL,
    height BIGINT NOT NULL,
    PRIMARY KEY (coins, height)
);
