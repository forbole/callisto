/* ---- PARAMS ---- */
CREATE TABLE staking_params
(
    bond_denom TEXT NOT NULL
);

/* ---- POOL ---- */

CREATE TABLE staking_pool_history
(
    bonded_tokens     BIGINT                      NOT NULL,
    not_bonded_tokens BIGINT                      NOT NULL,
    height            BIGINT                      NOT NULL,
    timestamp         TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    PRIMARY KEY (height)
);

/* ---- VALIDATORS INFO ---- */

CREATE TABLE validator_info
(
    consensus_address     TEXT NOT NULL UNIQUE REFERENCES validator (consensus_address),
    operator_address      TEXT NOT NULL UNIQUE,
    self_delegate_address TEXT REFERENCES account (address),
    max_change_rate       TEXT NOT NULL,
    max_rate              TEXT NOT NULL
);

CREATE TABLE validator_description
(
    validator_address TEXT NOT NULL UNIQUE REFERENCES validator (consensus_address),
    moniker           TEXT,
    identity          TEXT,
    website           TEXT,
    security_contact  TEXT,
    details           TEXT
);

CREATE TABLE validator_description_history
(
    validator_address TEXT                        NOT NULL REFERENCES validator (consensus_address),
    moniker           TEXT,
    identity          TEXT,
    website           TEXT,
    security_contact  TEXT,
    details           TEXT,
    height            BIGINT,
    timestamp         TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    PRIMARY KEY (validator_address, height)
);

CREATE TABLE validator_commission
(
    validator_address   TEXT    NOT NULL UNIQUE REFERENCES validator (consensus_address),
    commission          DECIMAL NOT NULL,
    min_self_delegation BIGINT  NOT NULL
);

CREATE TABLE validator_commission_history
(
    validator_address   TEXT                        NOT NULL REFERENCES validator (consensus_address),
    commission          DECIMAL                     NOT NULL,
    min_self_delegation BIGINT                      NOT NULL,
    height              BIGINT                      NOT NULL,
    timestamp           TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    PRIMARY KEY (validator_address, height)
);

CREATE TABLE validator_voting_power
(
    validator_address TEXT   NOT NULL UNIQUE REFERENCES validator (consensus_address),
    voting_power      BIGINT NOT NULL
);

CREATE TABLE validator_voting_power_history
(
    validator_address TEXT                        NOT NULL REFERENCES validator (consensus_address),
    voting_power      BIGINT                      NOT NULL,
    height            BIGINT                      NOT NULL,
    timestamp         TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    PRIMARY KEY (validator_address, height)
);

CREATE TABLE validator_uptime
(
    validator_address     TEXT   NOT NULL UNIQUE REFERENCES validator (consensus_address),
    signed_blocks_window  BIGINT NOT NULL,
    missed_blocks_counter BIGINT NOT NULL
);

CREATE TABLE validator_uptime_history
(
    validator_address     TEXT                        NOT NULL REFERENCES validator (consensus_address),
    signed_blocks_window  BIGINT                      NOT NULL,
    missed_blocks_counter BIGINT                      NOT NULL,
    height                BIGINT                      NOT NULL,
    timestamp             TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    PRIMARY KEY (validator_address, height)
);

CREATE TABLE validator_status
(
    validator_address TEXT    NOT NULL UNIQUE REFERENCES validator (consensus_address),
    status            INT     NOT NULL,
    jailed            BOOLEAN NOT NULL
);

CREATE TABLE validator_status_history
(
    validator_address TEXT                        NOT NULL REFERENCES validator (consensus_address),
    status            INT                         NOT NULL,
    jailed            BOOLEAN                     NOT NULL,
    height            BIGINT,
    timestamp         TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    PRIMARY KEY (validator_address, height)
);

/* ---- DELEGATIONS ---- */

/*
 * This table holds only the CURRENT delegations.
 * It should be updated on a BLOCK basis, deleting all the
 * existing data
 */
CREATE TABLE delegation
(
    validator_address TEXT    NOT NULL REFERENCES validator (consensus_address),
    delegator_address TEXT    NOT NULL REFERENCES account (address),
    amount            COIN    NOT NULL,
    shares            NUMERIC NOT NUll
);

/*
 * This table holds the HISTORICAL delegations.
 * It should be updated on a MESSAGE basis, to avoid data duplication.
 */
CREATE TABLE delegation_history
(
    validator_address TEXT                        NOT NULL REFERENCES validator (consensus_address),
    delegator_address TEXT                        NOT NULL REFERENCES account (address),
    amount            COIN                        NOT NULL,
    shares            NUMERIC                     NOT NUll,
    height            BIGINT                      NOT NULL,
    timestamp         TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

/*
 * This table holds only the CURRENT unbonding delegations.
 * It should be updated on a BLOCK basis, deleting all the
 * existing data
 */
CREATE TABLE unbonding_delegation
(
    validator_address    TEXT                        NOT NULL REFERENCES validator (consensus_address),
    delegator_address    TEXT                        NOT NULL REFERENCES account (address),
    amount               COIN                        NOT NUll,
    completion_timestamp TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

/*
 * This table holds the HISTORICAL unbonding delegations.
 * It should be updated on a MESSAGE basis, to avoid data duplication.
 */
CREATE TABLE unbonding_delegation_history
(
    validator_address    TEXT                        NOT NULL REFERENCES validator (consensus_address),
    delegator_address    TEXT                        NOT NULL REFERENCES account (address),
    amount               COIN                        NOT NUll,
    completion_timestamp TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    height               BIGINT                      NOT NULL,
    timestamp            TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

/*
 * This table holds only the CURRENT redelegations.
 * It should be updated on a BLOCK basis, deleting all the
 * existing data
 */
CREATE TABLE redelegation
(
    delegator_address     TEXT                        NOT NULL REFERENCES account (address),
    src_validator_address TEXT                        NOT NULL REFERENCES validator (consensus_address),
    dst_validator_address TEXT                        NOT NULL REFERENCES validator (consensus_address),
    amount                COIN                        NOT NULL,
    completion_time       TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

/*
 * This table holds the HISTORICAL redelegations.
 * It should be updated on a MESSAGE basis, to avoid data duplication.
 */
CREATE TABLE redelegation_history
(
    delegator_address     TEXT                        NOT NULL REFERENCES account (address),
    src_validator_address TEXT                        NOT NULL REFERENCES validator (consensus_address),
    dst_validator_address TEXT                        NOT NULL REFERENCES validator (consensus_address),
    amount                COIN                        NOT NULL,
    completion_time       TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    height                BIGINT                      NOT NULL,
    timestamp             TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

/*--------------------------------------------*/

/*
* This holds the votes that is the evidence of a double sign. This update on BLOCK basis when there is a double sign happens
*/
CREATE TABLE double_sign_vote
(
    signiture TEXT NOT NULL UNIQUE,
    hx TEXT NOT NULL,
    part_header TEXT NOT NULL,
    height                BIGINT                      NOT NULL,
    timestamp             TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

/**
* This holds the HISTORICAL double_sign_evidence on BLOCK basis
*/
CREATE TABLE double_sign_evidence
(
    pubkey TEXT NOT NULL,
    consensus_address     TEXT NOT NULL REFERENCES validator (consensus_address),
    vote_a TEXT NOT NULL REFERENCES double_sign_vote(signiture),
    vote_b text NOT NULL REFERENCES double_sign_vote(signiture),
    height                BIGINT                      NOT NULL,
    timestamp             TIMESTAMP WITHOUT TIME ZONE NOT NULL
);
