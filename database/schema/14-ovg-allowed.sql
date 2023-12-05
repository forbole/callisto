-- +migrate Up
CREATE TABLE IF NOT EXISTS overgold_allowed_addresses
(
    id               BIGSERIAL  NOT NULL PRIMARY KEY,
    creator          TEXT       NOT NULL,
    address          TEXT[]     NOT NULL
);

CREATE TABLE IF NOT EXISTS overgold_allowed_create_addresses
(
    id               BIGSERIAL  NOT NULL PRIMARY KEY,
    tx_hash          TEXT       NOT NULL,
    creator          TEXT       NOT NULL,
    address          TEXT[]     NOT NULL
);

CREATE TABLE IF NOT EXISTS overgold_allowed_delete_by_addresses
(
    id               BIGSERIAL  NOT NULL PRIMARY KEY,
    tx_hash          TEXT       NOT NULL,
    creator          TEXT       NOT NULL,
    address          TEXT[]     NOT NULL
);

CREATE TABLE IF NOT EXISTS overgold_allowed_delete_by_id
(
    id               BIGSERIAL  NOT NULL,
    tx_hash          TEXT       NOT NULL,
    creator          TEXT       NOT NULL
);

CREATE TABLE IF NOT EXISTS overgold_allowed_update_addresses
(
    id               BIGSERIAL  NOT NULL PRIMARY KEY,
    tx_hash          TEXT       NOT NULL,
    creator          TEXT       NOT NULL,
    address          TEXT[]     NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS overgold_allowed_update_addresses CASCADE;
DROP TABLE IF EXISTS overgold_allowed_delete_by_id CASCADE;
DROP TABLE IF EXISTS overgold_allowed_delete_by_addresses CASCADE;
DROP TABLE IF EXISTS overgold_allowed_create_addresses CASCADE;
DROP TABLE IF EXISTS overgold_allowed_addresses CASCADE;