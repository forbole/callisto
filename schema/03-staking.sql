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
    consensus_address     TEXT NOT NULL REFERENCES validator (consensus_address) UNIQUE PRIMARY KEY,
    operator_address      TEXT NOT NULL UNIQUE,
    self_delegate_address TEXT REFERENCES account (address)
);

CREATE TABLE validator_delegation
(
    consensus_address TEXT                        NOT NULL REFERENCES validator (consensus_address),
    delegator_address TEXT                        NOT NULL REFERENCES account (address),
    amount            COIN                        NOT NULL,
    height            BIGINT                      NOT NULL,
    timestamp         TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE TABLE validator_unbonding_delegation
(
    consensus_address    TEXT                        NOT NULL REFERENCES validator (consensus_address),
    delegator_address    TEXT                        NOT NULL REFERENCES account (address),
    amount               COIN                        NOT NUll,
    completion_timestamp TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    height               BIGINT                      NOT NULL,
    timestamp            TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE TABLE validator_redelegation
(
    delegator_address     TEXT                        NOT NULL REFERENCES account (address),
    src_validator_address TEXT                        NOT NULL REFERENCES validator (consensus_address),
    dst_validator_address TEXT                        NOT NULL REFERENCES validator (consensus_address),
    amount                COIN                        NOT NULL,
    height                BIGINT                      NOT NULL,
    completion_time       TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE TABLE validator_delegation_shares
(
    operator_address  TEXT                        NOT NULL REFERENCES validator_info (operator_address),
    delegator_address TEXT                        NOT NULL REFERENCES account (address),
    shares            NUMERIC                     NOT NUll,
    height            BIGINT                      NOT NULL,
    timestamp         TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    PRIMARY KEY (operator_address, delegator_address, height)
);

CREATE TABLE validator_commission
(
    operator_address    TEXT                        NOT NULL REFERENCES validator_info (operator_address) UNIQUE,
    timestamp           TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    commission          DECIMAL                     NOT NULL,
    min_self_delegation BIGINT                      NOT NULL,
    height              BIGINT                      NOT NULL
);

CREATE TABLE validator_voting_power
(
    consensus_address TEXT   NOT NULL REFERENCES validator (consensus_address),
    voting_power      BIGINT NOT NULL,
    height            BIGINT NOT NULL,
    PRIMARY KEY (consensus_address, height)
);

CREATE TABLE validator_description
(
    operator_address      TEXT NOT NULL REFERENCES validator_info(operator_address),
    moniker               TEXT ,
    identity              TEXT,
    website               TEXT,
    security_contact      TEXT,
    details               TEXT,
    height                BIGINT,
    timestamp             TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
