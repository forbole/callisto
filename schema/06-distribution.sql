CREATE TYPE DEC_COIN AS
(
    denom  TEXT,
    amount TEXT
);

CREATE TABLE community_pool
(
    coins  DEC_COIN[] NOT NULL,
    height BIGINT     NOT NULL,
    PRIMARY KEY (coins, height)
);
CREATE INDEX community_pool_height_index ON community_pool (height);

CREATE TABLE validator_commission_amount
(
    validator_address TEXT       NOT NULL REFERENCES validator (consensus_address),
    amount            DEC_COIN[] NOT NULL,
    height            BIGINT     NOT NULL REFERENCES block (height),
    UNIQUE (validator_address, height)
);
CREATE INDEX validator_commission_amount_validator_address_index ON validator_commission_amount (validator_address);
CREATE INDEX validator_commission_amount_height_index ON validator_commission_amount (height);

CREATE TABLE delegation_reward
(
    validator_address TEXT       NOT NULL REFERENCES validator (consensus_address),
    delegator_address TEXT       NOT NULL REFERENCES account (address),
    withdraw_address  TEXT       NOT NULL,
    amount            DEC_COIN[] NOT NULL,
    height            BIGINT     NOT NULL REFERENCES block (height),
    UNIQUE (validator_address, delegator_address, height)
);
CREATE INDEX delegation_reward_validator_address_index ON delegation_reward (validator_address);
CREATE INDEX delegation_reward_delegator_address_index ON delegation_reward (delegator_address);
CREATE INDEX delegation_reward_height_index ON delegation_reward (height);