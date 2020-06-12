CREATE TABLE validator
(
    id                SERIAL PRIMARY KEY,
    consensus_address CHARACTER VARYING(52) NOT NULL UNIQUE, /* Validator consensus address */
    consensus_pubkey  CHARACTER VARYING(83) NOT NULL UNIQUE
);

CREATE TABLE pre_commit
(
    id                SERIAL PRIMARY KEY,
    validator_address CHARACTER VARYING(52)       NOT NULL REFERENCES validator (consensus_address),
    timestamp         TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    voting_power      INTEGER                     NOT NULL,
    proposer_priority INTEGER                     NOT NULL
);

CREATE TABLE block
(
    id               SERIAL PRIMARY KEY,
    height           INTEGER                     NOT NULL UNIQUE,
    hash             CHARACTER VARYING(64)       NOT NULL UNIQUE,
    num_txs          INTEGER DEFAULT 0,
    total_gas        INTEGER DEFAULT 0,
    proposer_address CHARACTER VARYING(52)       NOT NULL REFERENCES validator (consensus_address),
    pre_commits      INTEGER                     NOT NULL,
    timestamp        TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE TABLE transaction
(
    id         SERIAL PRIMARY KEY,
    timestamp  TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    gas_wanted INTEGER                              DEFAULT 0,
    gas_used   INTEGER                              DEFAULT 0,
    height     INTEGER                     NOT NULL REFERENCES block (height),
    txhash     CHARACTER VARYING(64)       NOT NULL UNIQUE,
    messages   JSONB                       NOT NULL DEFAULT '[]'::JSONB,
    fee        JSONB                       NOT NULL DEFAULT '{}'::JSONB,
    signatures JSONB                       NOT NULL DEFAULT '[]'::JSONB,
    memo       CHARACTER VARYING(256)
);
