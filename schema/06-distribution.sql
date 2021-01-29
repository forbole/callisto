CREATE TYPE DEC_COIN AS
(
    denom  TEXT,
    amount NUMERIC
);

CREATE TABLE community_pool
(
    coins  DEC_COIN[] NOT NULL,
    height BIGINT     NOT NULL,
    PRIMARY KEY (coins, height)
);

CREATE TABLE validator_commission_amount
(
    validator_address TEXT       NOT NULL REFERENCES validator (consensus_address),
    amount            DEC_COIN[] NOT NULL,
    height            BIGINT     NOT NULL REFERENCES block (height),
    UNIQUE (validator_address, height)
);

CREATE INDEX validator_commission_amount_height_index ON validator_commission_amount (height);

CREATE TABLE delegation_amount
(
    validator_address TEXT       NOT NULL REFERENCES validator (consensus_address),
    delegator_address TEXT       NOT NULL REFERENCES account (address),
    amount            DEC_COIN[] NOT NULL,
    height            BIGINT     NOT NULL REFERENCES block (height),
    UNIQUE (validator_address, delegator_address, height)
);

CREATE INDEX delegator_reward_amount_height_index ON delegation_amount (height);