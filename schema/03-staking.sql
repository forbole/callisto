/* ---- PARAMS ---- */
CREATE TABLE staking_params
(
    bond_denom TEXT NOT NULL
);

/* ---- POOL ---- */

CREATE TABLE staking_pool_history
(
    bonded_tokens     BIGINT NOT NULL,
    not_bonded_tokens BIGINT NOT NULL,
    height            BIGINT NOT NULL,
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
    validator_address TEXT NOT NULL REFERENCES validator (consensus_address),
    moniker           TEXT,
    identity          TEXT,
    website           TEXT,
    security_contact  TEXT,
    details           TEXT,
    height            BIGINT,
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
    validator_address   TEXT    NOT NULL REFERENCES validator (consensus_address),
    commission          DECIMAL NOT NULL,
    min_self_delegation BIGINT  NOT NULL,
    height              BIGINT  NOT NULL,
    PRIMARY KEY (validator_address, height)
);

CREATE TABLE validator_voting_power
(
    validator_address TEXT   NOT NULL UNIQUE REFERENCES validator (consensus_address),
    voting_power      BIGINT NOT NULL
);

CREATE TABLE validator_voting_power_history
(
    validator_address TEXT   NOT NULL REFERENCES validator (consensus_address),
    voting_power      BIGINT NOT NULL,
    height            BIGINT NOT NULL REFERENCES block (height),
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
    validator_address     TEXT   NOT NULL REFERENCES validator (consensus_address),
    signed_blocks_window  BIGINT NOT NULL,
    missed_blocks_counter BIGINT NOT NULL,
    height                BIGINT NOT NULL,
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
    validator_address TEXT    NOT NULL REFERENCES validator (consensus_address),
    status            INT     NOT NULL,
    jailed            BOOLEAN NOT NULL,
    height            BIGINT,
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
    id                SERIAL  NOT NULL PRIMARY KEY,
    validator_address TEXT    NOT NULL REFERENCES validator (consensus_address),
    delegator_address TEXT    NOT NULL REFERENCES account (address),
    amount            COIN    NOT NULL,
    shares            NUMERIC NOT NUll
);

/**
  * This function is used to have a Hasura compute field (https://hasura.io/docs/1.0/graphql/core/schema/computed-fields.html)
  * inside the delegation table, so that it's easy to determine whether an entry represents a self delegation or not.
 */
CREATE FUNCTION is_delegation_self_delegate(delegation_row delegation)
    RETURNS BOOLEAN AS
$$
SELECT (
           SELECT self_delegate_address
           FROM validator_info
           WHERE validator_info.consensus_address = delegation_row.validator_address
       ) = delegation_row.delegator_address
$$ LANGUAGE sql STABLE;

/**
  * This function is used to add a self_delegations field to the validator table allowing to easily get all the
  * self delegations related to a specific validator.
 */
CREATE FUNCTION self_delegations(validator_row validator) RETURNS SETOF delegation AS
$$
SELECT *
FROM delegation
WHERE delegator_address = (
    SELECT self_delegate_address
    FROM validator_info
    WHERE validator_info.consensus_address = validator_row.consensus_address
)
$$
    LANGUAGE sql
    STABLE;

/*
 * This table holds the HISTORICAL delegations.
 * It should be updated on a MESSAGE basis, to avoid data duplication.
 */
CREATE TABLE delegation_history
(
    id                SERIAL  NOT NULL PRIMARY KEY,
    validator_address TEXT    NOT NULL REFERENCES validator (consensus_address),
    delegator_address TEXT    NOT NULL REFERENCES account (address),
    amount            COIN    NOT NULL,
    shares            NUMERIC NOT NUll,
    height            BIGINT  NOT NULL
);

/**
  * This function is used to have a Hasura compute field (https://hasura.io/docs/1.0/graphql/core/schema/computed-fields.html)
  * inside the delegation_history table, so that it's easy to determine whether an entry represents a self delegation or not.
 */
CREATE FUNCTION is_delegation_history_self_delegate(delegation_row delegation_history)
    RETURNS BOOLEAN AS
$$
SELECT (
           SELECT self_delegate_address
           FROM validator_info
           WHERE validator_info.consensus_address = delegation_row.validator_address
       ) = delegation_row.delegator_address
$$ LANGUAGE sql STABLE;

/**
  * This function is used to add a self_delegation_histories field to the validator table allowing to easily get all the
  * self delegation histories related to a specific validator.
 */
CREATE FUNCTION self_delegation_histories(validator_row validator) RETURNS SETOF delegation_history AS
$$
SELECT *
FROM delegation_history
WHERE delegation_history.delegator_address = (
    SELECT self_delegate_address
    FROM validator_info
    WHERE validator_info.consensus_address = validator_row.consensus_address
)
$$
    LANGUAGE sql
    STABLE;

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
    height               BIGINT                      NOT NULL
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
    height                BIGINT                      NOT NULL
);

/*--------------------------------------------*/

/*
 * This holds the votes that is the evidence of a double sign.
 * It should be updated on a BLOCK basis when a double sign occurs.
 */
CREATE TABLE double_sign_vote
(
    id                SERIAL PRIMARY KEY,
    type              SMALLINT NOT NULL,
    height            BIGINT   NOT NULL,
    round             INT      NOT NULL,
    block_id          TEXT     NOT NULL,
    validator_address TEXT     NOT NULL REFERENCES validator (consensus_address),
    validator_index   INT      NOT NULL,
    signature         TEXT     NOT NULL,
    UNIQUE (block_id, validator_address)
);

/*
 * This holds the HISTORICAL double_sign_evidence.
 * It should be updated on a on BLOCK basis.
 */
CREATE TABLE double_sign_evidence
(
    public_key TEXT   NOT NULL,
    vote_a_id  BIGINT NOT NULL REFERENCES double_sign_vote (id),
    vote_b_id  BIGINT NOT NULL REFERENCES double_sign_vote (id)
);
