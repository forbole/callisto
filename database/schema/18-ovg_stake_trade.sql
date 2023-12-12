-- +migrate Up
CREATE TABLE IF NOT EXISTS overgold_stake_sell
(
    id               BIGSERIAL  NOT NULL PRIMARY KEY,
    tx_hash          TEXT       NOT NULL,
    creator          TEXT       NOT NULL,
    amount           BIGINT     NOT NULL
);

CREATE TABLE IF NOT EXISTS overgold_stake_buy
(
    id               BIGSERIAL  NOT NULL PRIMARY KEY,
    tx_hash          TEXT       NOT NULL,
    creator          TEXT       NOT NULL,
    amount           BIGINT     NOT NULL
);

CREATE TABLE IF NOT EXISTS overgold_stake_sell_cancel
(
    id               BIGSERIAL  NOT NULL PRIMARY KEY,
    tx_hash          TEXT       NOT NULL,
    creator          TEXT       NOT NULL,
    amount           BIGINT     NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS overgold_stake_sell CASCADE;
DROP TABLE IF EXISTS overgold_stake_buy CASCADE;
DROP TABLE IF EXISTS overgold_stake_sell_cancel CASCADE;