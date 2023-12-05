-- +migrate Up
CREATE TABLE IF NOT EXISTS overgold_stake_distribute_rewards
(
    id      BIGSERIAL NOT NULL PRIMARY KEY,
    tx_hash TEXT      NOT NULL,
    creator TEXT      NOT NULL
);

CREATE TABLE IF NOT EXISTS overgold_stake_claim_reward
(
    id      BIGSERIAL NOT NULL PRIMARY KEY,
    tx_hash TEXT      NOT NULL,
    creator TEXT      NOT NULL,
    amount  BIGINT    NOT NULL
);

-- +migrate Down
DROP TABLE IF EXISTS overgold_stake_distribute_rewards CASCADE;
DROP TABLE IF EXISTS overgold_stake_claim_reward CASCADE;