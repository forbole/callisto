CREATE TABLE validator
(
    consensus_address TEXT NOT NULL UNIQUE PRIMARY KEY,
    consensus_pubkey  TEXT NOT NULL UNIQUE
);

CREATE TABLE pre_commit
(
    id                SERIAL PRIMARY KEY,
    height            BIGINT                      NOT NULL,
    validator_address TEXT                        NOT NULL REFERENCES validator (consensus_address),
    timestamp         TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    voting_power      BIGINT                      NOT NULL,
    proposer_priority INTEGER                     NOT NULL
);

CREATE TABLE block
(
    height           BIGINT                      NOT NULL UNIQUE PRIMARY KEY,
    hash             TEXT                        NOT NULL UNIQUE,
    num_txs          INTEGER DEFAULT 0,
    total_gas        INTEGER DEFAULT 0,
    proposer_address TEXT                        NOT NULL REFERENCES validator (consensus_address),
    pre_commits      INTEGER                     NOT NULL,
    timestamp        TIMESTAMP WITHOUT TIME ZONE NOT NULL
);

CREATE TABLE transaction
(
    hash         TEXT    NOT NULL UNIQUE PRIMARY KEY,
    height       BIGINT  NOT NULL REFERENCES block (height),
    success      BOOLEAN NOT NULL,
    messages     JSONB   NOT NULL DEFAULT '[]'::JSONB,
    memo         TEXT,
    signatures   TEXT[]  NOT NULL,
    signer_infos JSONB   NOT NULL DEFAULT '[]'::JSONB,
    fee          JSONB   NOT NULL DEFAULT '{}'::JSONB,
    gas_wanted   INTEGER          DEFAULT 0,
    gas_used     INTEGER          DEFAULT 0,
    raw_log      TEXT,
    logs         JSONB
);

