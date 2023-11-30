-- +migrate Up
CREATE TYPE SEND_DATA AS
(
    address          TEXT,
    coins            COIN[]
);

CREATE TABLE IF NOT EXISTS msg_multi_send
(
    id               BIGSERIAL      NOT NULL PRIMARY KEY,
    tx_hash          TEXT           NOT NULL,
    inputs           COIN[]         NOT NULL,
    outputs          COIN[]         NOT NULL
);

CREATE TABLE IF NOT EXISTS msg_send
(
    id               BIGSERIAL      NOT NULL PRIMARY KEY,
    tx_hash          TEXT           NOT NULL,
    from_address     TEXT           NOT NULL,
    to_address       TEXT           NOT NULL,
    amount           COIN[]         NOT NULL DEFAULT '{}'
);

-- +migrate Down
DROP TABLE IF EXISTS msg_send CASCADE;
DROP TABLE IF EXISTS msg_multi_send CASCADE;
DROP TYPE IF EXISTS SEND_DATA CASCADE;