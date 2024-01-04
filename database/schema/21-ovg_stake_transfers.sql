-- +migrate Up
CREATE TABLE IF NOT EXISTS overgold_stake_transfer_from_user
(
    id               BIGSERIAL  NOT NULL PRIMARY KEY,
    tx_hash          TEXT       NOT NULL,
    creator          TEXT       NOT NULL,
    amount           BIGINT     NOT NULL,
    address          TEXT       NOT NULL
);

CREATE TABLE IF NOT EXISTS overgold_stake_transfer_to_user
(
    id               BIGSERIAL  NOT NULL PRIMARY KEY,
    tx_hash          TEXT       NOT NULL,
    creator          TEXT       NOT NULL,
    amount           BIGINT     NOT NULL,
    address          TEXT       NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS overgold_stake_transfer_from_user CASCADE;
DROP TABLE IF EXISTS overgold_stake_transfer_to_user CASCADE;
