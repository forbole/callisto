CREATE TABLE staking_pool
(
    height            BIGINT                      NOT NULL PRIMARY KEY,
    timestamp         TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    bonded_tokens     BIGINT                      NOT NULL,
    not_bonded_tokens BIGINT                      NOT NULL
);

CREATE TABLE validator_uptime
(
    validator_address     TEXT   NOT NULL REFERENCES validator (consensus_address),
    height                BIGINT NOT NULL,
    signed_blocks_window  BIGINT NOT NULL,
    missed_blocks_counter BIGINT NOT NULL,
    PRIMARY KEY (validator_address, height)
);

CREATE TABLE validator_info
(
    consensus_address TEXT NOT NULL REFERENCES validator (consensus_address) UNIQUE PRIMARY KEY,
    operator_address  TEXT NOT NULL
);

CREATE TABLE validator_delegations
(
    consensus_address TEXT                        NOT NULL REFERENCES validator (consensus_address),
    delegator_address TEXT                        NOT NULL REFERENCES account (address),
    shares            DECIMAL                     NOT NULL,
    balance           COIN,
    height            BIGINT                      NOT NULL,
    timestamp         TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    PRIMARY KEY (consensus_address, delegator_address, height)
);

CREATE TABLE validator_unbonding_delegations
(
    consensus_address TEXT                        NOT NULL REFERENCES validator (consensus_address),
    delegator_address TEXT                        NOT NULL REFERENCES account (address),
    initial_balance   DECIMAL                     NOT NUll,
    balance           DECIMAL                     NOT NULL,
    creation_height   BIGINT                      NOT NULL,
    height            BIGINT                      NOT NULL,
    timestamp         TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    PRIMARY KEY (consensus_address, delegator_address, height)
)
