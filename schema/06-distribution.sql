CREATE TYPE DEC_COIN AS
(
    denom  TEXT,
    amount TEXT
);

CREATE TABLE community_pool
(
    one_row_id bool PRIMARY KEY DEFAULT TRUE,
    coins      DEC_COIN[] NOT NULL,
    height     BIGINT     NOT NULL,
    CONSTRAINT one_row_uni CHECK (one_row_id)
);
CREATE INDEX community_pool_height_index ON community_pool (height);

CREATE TABLE validator_commission_amount
(
    validator_address TEXT       NOT NULL REFERENCES validator (consensus_address) PRIMARY KEY,
    amount            DEC_COIN[] NOT NULL,
    height            BIGINT     NOT NULL
);
CREATE INDEX validator_commission_amount_height_index ON validator_commission_amount (height);

CREATE TABLE delegation_reward
(
    validator_address TEXT       NOT NULL REFERENCES validator (consensus_address),
    delegator_address TEXT       NOT NULL REFERENCES account (address),
    withdraw_address  TEXT       NOT NULL,
    amount            DEC_COIN[] NOT NULL,
    height            BIGINT     NOT NULL,
    CONSTRAINT validator_delegator_unique UNIQUE (validator_address, delegator_address)
);
CREATE INDEX delegation_reward_delegator_address_index ON delegation_reward (delegator_address);
CREATE INDEX delegation_reward_height_index ON delegation_reward (height);