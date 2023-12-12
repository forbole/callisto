-- +migrate Up
CREATE TABLE IF NOT EXISTS overgold_core_issue
(
    id               BIGSERIAL  NOT NULL PRIMARY KEY,
    tx_hash          TEXT       NOT NULL,
    creator          TEXT       NOT NULL,
    amount           BIGINT     NOT NULL,
    denom            TEXT       NOT NULL,
    address          TEXT       NOT NULL
);

CREATE TABLE IF NOT EXISTS overgold_core_withdraw
(
    id               BIGSERIAL  NOT NULL PRIMARY KEY,
    tx_hash          TEXT       NOT NULL,
    creator          TEXT       NOT NULL,
    amount           BIGINT     NOT NULL,
    denom            TEXT       NOT NULL,
    address          TEXT       NOT NULL
);

CREATE TABLE IF NOT EXISTS overgold_core_send
(
    id               BIGSERIAL  NOT NULL PRIMARY KEY,
    tx_hash          TEXT       NOT NULL,
    creator          TEXT       NOT NULL,
    address_from     TEXT       NOT NULL,
    address_to       TEXT       NOT NULL,
    amount           BIGINT     NOT NULL,
    denom            TEXT       NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS overgold_core_send CASCADE;
DROP TABLE IF EXISTS overgold_core_withdraw CASCADE;
DROP TABLE IF EXISTS overgold_core_issue CASCADE;