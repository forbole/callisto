CREATE TABLE validator
(
    consensus_address TEXT NOT NULL UNIQUE PRIMARY KEY,
    consensus_pubkey  TEXT NOT NULL UNIQUE
);

CREATE TABLE pre_commit
(
    id                SERIAL PRIMARY KEY,
    validator_address TEXT                        NOT NULL REFERENCES validator (consensus_address),
    timestamp         TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    voting_power      INTEGER                     NOT NULL,
    proposer_priority INTEGER                     NOT NULL
);

CREATE TABLE block
(
    height           INTEGER                     NOT NULL UNIQUE PRIMARY KEY,
    hash             TEXT                        NOT NULL UNIQUE,
    num_txs          INTEGER DEFAULT 0,
    total_gas        INTEGER DEFAULT 0,
    proposer_address TEXT                        NOT NULL REFERENCES validator (consensus_address),
    pre_commits      INTEGER                     NOT NULL,
    timestamp        TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE TABLE transaction
(
    timestamp  TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    gas_wanted INTEGER                              DEFAULT 0,
    gas_used   INTEGER                              DEFAULT 0,
    height     INTEGER                     NOT NULL REFERENCES block (height),
    hash       TEXT                        NOT NULL UNIQUE PRIMARY KEY,
    messages   JSONB                       NOT NULL DEFAULT '[]'::JSONB,
    fee        JSONB                       NOT NULL DEFAULT '{}'::JSONB,
    signatures JSONB                       NOT NULL DEFAULT '[]'::JSONB,
    memo       TEXT
);

