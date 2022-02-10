/* ---- PARAMS ---- */

CREATE TABLE staking_params
(
    one_row_id BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    params     JSONB   NOT NULL,
    height     BIGINT  NOT NULL,
    CHECK (one_row_id)
);
CREATE INDEX staking_params_height_index ON staking_params (height);

/* ---- POOL ---- */

CREATE TABLE staking_pool
(
    one_row_id        BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    bonded_tokens     TEXT    NOT NULL,
    not_bonded_tokens TEXT    NOT NULL,
    height            BIGINT  NOT NULL,
    CHECK (one_row_id)
);
CREATE INDEX staking_pool_height_index ON staking_pool (height);

/* ---- VALIDATORS INFO ---- */

CREATE TABLE validator_info
(
    consensus_address     TEXT   NOT NULL UNIQUE PRIMARY KEY REFERENCES validator (consensus_address),
    operator_address      TEXT   NOT NULL UNIQUE,
    self_delegate_address TEXT REFERENCES account (address),
    max_change_rate       TEXT   NOT NULL,
    max_rate              TEXT   NOT NULL,
    height                BIGINT NOT NULL
);
CREATE INDEX validator_info_operator_address_index ON validator_info (operator_address);
CREATE INDEX validator_info_self_delegate_address_index ON validator_info (self_delegate_address);

CREATE TABLE validator_description
(
    validator_address TEXT   NOT NULL REFERENCES validator (consensus_address) PRIMARY KEY,
    moniker           TEXT,
    identity          TEXT,
    avatar_url        TEXT,
    website           TEXT,
    security_contact  TEXT,
    details           TEXT,
    height            BIGINT NOT NULL
);
CREATE INDEX validator_description_height_index ON validator_description (height);

CREATE TABLE validator_commission
(
    validator_address   TEXT    NOT NULL REFERENCES validator (consensus_address) PRIMARY KEY,
    commission          DECIMAL NOT NULL,
    min_self_delegation BIGINT  NOT NULL,
    height              BIGINT  NOT NULL
);
CREATE INDEX validator_commission_height_index ON validator_commission (height);

CREATE TABLE validator_voting_power
(
    validator_address TEXT   NOT NULL REFERENCES validator (consensus_address) PRIMARY KEY,
    voting_power      BIGINT NOT NULL,
    height            BIGINT NOT NULL REFERENCES block (height)
);
CREATE INDEX validator_voting_power_height_index ON validator_voting_power (height);

CREATE TABLE validator_status
(
    validator_address TEXT    NOT NULL REFERENCES validator (consensus_address) PRIMARY KEY,
    status            INT     NOT NULL,
    jailed            BOOLEAN NOT NULL,
    tombstoned        BOOLEAN NOT NULL DEFAULT FALSE,
    height            BIGINT  NOT NULL
);
CREATE INDEX validator_status_height_index ON validator_status (height);

/* ---- DOUBLE SIGN EVIDENCE ---- */

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
CREATE INDEX double_sign_vote_validator_address_index ON double_sign_vote (validator_address);
CREATE INDEX double_sign_vote_height_index ON double_sign_vote (height);

/*
 * This holds the double sign evidences.
 * It should be updated on a on BLOCK basis.
 */
CREATE TABLE double_sign_evidence
(
    height    BIGINT NOT NULL,
    vote_a_id BIGINT NOT NULL REFERENCES double_sign_vote (id),
    vote_b_id BIGINT NOT NULL REFERENCES double_sign_vote (id)
);
CREATE INDEX double_sign_evidence_height_index ON double_sign_evidence (height);