CREATE TYPE DEC_COIN AS
(
    denom  TEXT,
    amount TEXT
);

CREATE TABLE community_pool
(
    coins  DEC_COIN[] NOT NULL
);

CREATE TABLE validator_commission_amount
(
    validator_address TEXT       NOT NULL REFERENCES validator (consensus_address) PRIMARY KEY,
    amount            DEC_COIN[] NOT NULL
);
CREATE INDEX validator_commission_amount_validator_address_index ON validator_commission_amount (validator_address);

CREATE TABLE delegation_reward
(
    validator_address TEXT       NOT NULL REFERENCES validator (consensus_address),
    delegator_address TEXT       NOT NULL REFERENCES account (address),
    withdraw_address  TEXT       NOT NULL,
    amount            DEC_COIN[] NOT NULL,
    CONSTRAINT validator_delegator_unique UNIQUE (validator_address, delegator_address, withdraw_address)
);
CREATE INDEX delegation_reward_validator_address_index ON delegation_reward (validator_address);
CREATE INDEX delegation_reward_delegator_address_index ON delegation_reward (delegator_address);