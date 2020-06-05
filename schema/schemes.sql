--- COSMOS ----------------------------------------------

CREATE TABLE validator
(
    id               SERIAL PRIMARY KEY,
    address          character varying(40) NOT NULL UNIQUE,
    consensus_pubkey character varying(83) NOT NULL UNIQUE
);

CREATE TABLE pre_commit
(
    id                SERIAL PRIMARY KEY,
    validator_address character varying(40)       NOT NULL REFERENCES validator (address),
    timestamp         timestamp without time zone NOT NULL,
    voting_power      integer                     NOT NULL,
    proposer_priority integer                     NOT NULL
);

CREATE TABLE block
(
    id               SERIAL PRIMARY KEY,
    height           integer                     NOT NULL UNIQUE,
    hash             character varying(64)       NOT NULL UNIQUE,
    num_txs          integer DEFAULT 0,
    total_gas        integer DEFAULT 0,
    proposer_address character varying(40)       NOT NULL REFERENCES validator (address),
    pre_commits      integer                     NOT NULL,
    timestamp        timestamp without time zone NOT NULL
);

CREATE TABLE transaction
(
    id         SERIAL PRIMARY KEY,
    timestamp  timestamp without time zone NOT NULL,
    gas_wanted integer                              DEFAULT 0,
    gas_used   integer                              DEFAULT 0,
    height     integer                     NOT NULL REFERENCES block (height),
    txhash     character varying(64)       NOT NULL UNIQUE,
    messages   jsonb                       NOT NULL DEFAULT '[]'::jsonb,
    fee        jsonb                       NOT NULL DEFAULT '{}'::jsonb,
    signatures jsonb                       NOT NULL DEFAULT '[]'::jsonb,
    memo       character varying(256)
);
