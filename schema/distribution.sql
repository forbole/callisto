CREATE TABLE community_pool
(
    coins                  COIN[]               NOT NULL,
    height                BIGINT                NOT NULL,
    PRIMARY KEY (coins,height)
);