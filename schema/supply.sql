CREATE TABLE total_coins
(
    coins                  COIN[]               NOT NULL,
    height                BIGINT                NOT NULL,
    PRIMARY KEY (coins,height)
);