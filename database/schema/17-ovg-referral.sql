-- +migrate Up
CREATE TABLE IF NOT EXISTS overgold_referral_set_referrer
(
    id               BIGSERIAL  NOT NULL PRIMARY KEY,
    tx_hash          TEXT       NOT NULL,
    creator          TEXT       NOT NULL,
    referrer_address TEXT       NOT NULL,
    referral_address TEXT       NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS overgold_referral_set_referrer CASCADE;