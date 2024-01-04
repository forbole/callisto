-- +migrate Up
CREATE TABLE overgold_feeexcluder_address ( 
    id          BIGSERIAL   NOT NULL PRIMARY KEY,
    msg_id      BIGINT,
    address     TEXT        NOT NULL,
    creator     TEXT        NOT NULL
);

CREATE TABLE overgold_feeexcluder_daily_stats ( 
    id              BIGSERIAL   NOT NULL PRIMARY KEY,
    msg_id          BIGINT      NOT NULL,
    amount_with_fee COIN[],
    amount_no_fee   COIN[],
    fee             COIN[],
    count_with_fee  INT         NOT NULL,
    count_no_fee    INT         NOT NULL
);

CREATE TABLE overgold_feeexcluder_stats ( 
    id              TEXT        NOT NULL PRIMARY KEY,
    date            TIMESTAMP   NOT NULL,
    daily_stats_id  BIGSERIAL   NOT NULL REFERENCES overgold_feeexcluder_daily_stats(id)
);

CREATE TABLE overgold_feeexcluder_fees ( 
    id              BIGSERIAL   NOT NULL PRIMARY KEY,
    msg_id          BIGINT,
    creator         TEXT        NOT NULL,
    amount_from     BIGSERIAL   NOT NULL,
    fee             NUMERIC     NOT NULL,
    ref_reward      NUMERIC     NOT NULL,
    stake_reward    NUMERIC     NOT NULL,
    min_amount      BIGSERIAL   NOT NULL,
    no_ref_reward   BOOLEAN     NOT NULL
);

CREATE TABLE overgold_feeexcluder_tariff (
    id              BIGSERIAL   NOT NULL PRIMARY KEY,
    msg_id          BIGINT,
    amount          BIGSERIAL   NOT NULL,
    denom           TEXT        NOT NULL,
    min_ref_balance BIGSERIAL   NOT NULL
);

CREATE TABLE overgold_feeexcluder_m2m_tariff_fees (
    tariff_id BIGSERIAL REFERENCES overgold_feeexcluder_tariff(id),
    fees_id   BIGSERIAL REFERENCES overgold_feeexcluder_fees(id),
    PRIMARY KEY (tariff_id, fees_id)
);

CREATE TABLE overgold_feeexcluder_tariffs (
    id          BIGSERIAL   NOT NULL PRIMARY KEY,
    denom       TEXT        NOT NULL,
    creator     TEXT        NOT NULL
);

CREATE TABLE overgold_feeexcluder_m2m_tariff_tariffs (
    tariff_id  BIGSERIAL REFERENCES overgold_feeexcluder_tariff(id),
    tariffs_id BIGSERIAL REFERENCES overgold_feeexcluder_tariffs(id),
    PRIMARY KEY (tariff_id, tariffs_id)
);

CREATE TABLE overgold_feeexcluder_genesis_state (
    id                  BIGSERIAL NOT NULL PRIMARY KEY,
    address_count       BIGSERIAL,
    daily_stats_count   BIGSERIAL
);

CREATE TABLE overgold_feeexcluder_m2m_genesis_state_address (
    genesis_state_id  BIGSERIAL REFERENCES overgold_feeexcluder_genesis_state(id),
    address_id        BIGSERIAL REFERENCES overgold_feeexcluder_address(id),
    PRIMARY KEY (genesis_state_id, address_id)
);

CREATE TABLE overgold_feeexcluder_m2m_genesis_state_daily_stats (
    genesis_state_id  BIGSERIAL REFERENCES overgold_feeexcluder_genesis_state(id),
    daily_stats_id    BIGSERIAL REFERENCES overgold_feeexcluder_daily_stats(id),
    PRIMARY KEY (genesis_state_id, daily_stats_id)
);

CREATE TABLE overgold_feeexcluder_m2m_genesis_state_stats (
    genesis_state_id  BIGSERIAL REFERENCES overgold_feeexcluder_genesis_state(id),
    stats_id          TEXT      REFERENCES overgold_feeexcluder_stats(id),
    PRIMARY KEY (genesis_state_id, stats_id)
);

CREATE TABLE overgold_feeexcluder_m2m_genesis_state_tariffs (
    genesis_state_id  BIGSERIAL REFERENCES overgold_feeexcluder_genesis_state(id),
    tariffs_id        BIGINT    REFERENCES overgold_feeexcluder_tariffs(id),
    PRIMARY KEY (genesis_state_id, tariffs_id)
);

CREATE TABLE IF NOT EXISTS overgold_feeexcluder_create_address
(
    id               BIGSERIAL  NOT NULL PRIMARY KEY,
    tx_hash          TEXT       NOT NULL,
    creator          TEXT       NOT NULL,
    address          TEXT       NOT NULL
);

CREATE TABLE IF NOT EXISTS overgold_feeexcluder_update_address
(
    id               BIGSERIAL  NOT NULL PRIMARY KEY,
    tx_hash          TEXT       NOT NULL,
    creator          TEXT       NOT NULL,
    address          TEXT       NOT NULL,
    denom            TEXT       NOT NULL
);

CREATE TABLE IF NOT EXISTS overgold_feeexcluder_delete_address
(
    id               BIGSERIAL  NOT NULL PRIMARY KEY,
    tx_hash          TEXT       NOT NULL,
    creator          TEXT       NOT NULL
);

CREATE TABLE IF NOT EXISTS overgold_feeexcluder_create_tariffs
(
    id               BIGSERIAL  NOT NULL PRIMARY KEY,
    tx_hash          TEXT       NOT NULL,
    creator          TEXT       NOT NULL,
    denom            TEXT       NOT NULL,
    tariff_id        BIGINT     NOT NULL REFERENCES overgold_feeexcluder_tariff(id)
);

CREATE TABLE IF NOT EXISTS overgold_feeexcluder_update_tariffs
(
    id               BIGSERIAL  NOT NULL PRIMARY KEY,
    tx_hash          TEXT       NOT NULL,
    creator          TEXT       NOT NULL,
    denom            TEXT       NOT NULL,
    tariff_id        BIGINT     NOT NULL REFERENCES overgold_feeexcluder_tariff(id)
);

CREATE TABLE IF NOT EXISTS overgold_feeexcluder_delete_tariffs
(
    id               BIGSERIAL  NOT NULL PRIMARY KEY,
    tx_hash          TEXT       NOT NULL,
    creator          TEXT       NOT NULL,
    denom            TEXT       NOT NULL,
    tariff_id        BIGINT     NOT NULL REFERENCES overgold_feeexcluder_tariff(id),
    fees_id          BIGINT     NOT NULL REFERENCES overgold_feeexcluder_fees(id)
);

-- +migrate Down
DROP TABLE IF EXISTS overgold_feeexcluder_delete_tariffs CASCADE;
DROP TABLE IF EXISTS overgold_feeexcluder_update_tariffs CASCADE;
DROP TABLE IF EXISTS overgold_feeexcluder_create_tariffs CASCADE;
DROP TABLE IF EXISTS overgold_feeexcluder_tariff CASCADE;
DROP TABLE IF EXISTS overgold_feeexcluder_delete_address CASCADE;
DROP TABLE IF EXISTS overgold_feeexcluder_update_address CASCADE;
DROP TABLE IF EXISTS overgold_feeexcluder_create_address CASCADE;
DROP TABLE IF EXISTS overgold_feeexcluder_m2m_genesis_state_tariffs CASCADE;
DROP TABLE IF EXISTS overgold_feeexcluder_m2m_genesis_state_stats CASCADE;
DROP TABLE IF EXISTS overgold_feeexcluder_m2m_genesis_state_daily_stats CASCADE;
DROP TABLE IF EXISTS overgold_feeexcluder_m2m_genesis_state_address CASCADE;
DROP TABLE IF EXISTS overgold_feeexcluder_genesis_state CASCADE;
DROP TABLE IF EXISTS overgold_feeexcluder_m2m_tariff_tariffs CASCADE;
DROP TABLE IF EXISTS overgold_feeexcluder_tariffs CASCADE;
DROP TABLE IF EXISTS overgold_feeexcluder_m2m_tariff_fees CASCADE;
DROP TABLE IF EXISTS overgold_feeexcluder_tariff CASCADE;
DROP TABLE IF EXISTS overgold_feeexcluder_fees CASCADE;
DROP TABLE IF EXISTS overgold_feeexcluder_stats CASCADE;
DROP TABLE IF EXISTS overgold_feeexcluder_daily_stats CASCADE;
DROP TABLE IF EXISTS overgold_feeexcluder_address CASCADE;