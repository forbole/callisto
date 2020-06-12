---------------------------------------------------------
--- STANDARD --------------------------------------------
---------------------------------------------------------

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

---------------------------------------------------------
--- STAKING ---------------------------------------------
---------------------------------------------------------

CREATE TABLE staking_pool
(
    height            BIGINT                      NOT NULL,
    timestamp         TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    bonded_tokens     BIGINT                      NOT NULL,
    not_bonded_tokens BIGINT                      NOT NULL
);

CREATE TABLE validator_description
(
    id               SERIAL PRIMARY KEY,
    moniker          TEXT,
    identify         TEXT,
    website          TEXT,
    security_contact TEXT,
    DETAILS          TEXT
);

CREATE TABLE validator_commission
(
    id              SERIAL PRIMARY KEY,
    current_rate    DECIMAL NOT NULL,
    max_rate        DECIMAL NOT NULL,
    max_rate_change DECIMAL NOT NULL,
    update_time     DATE
);

CREATE TABLE validator_staking_info
(
    validator_address      character varying(40) NOT NULL REFERENCES validator (address),
    operator_address       TEXT                  NOT NULL,
    consensus_pubkey       TEXT                  NOT NULL,
    jailed                 BOOLEAN               NOT NULL,
    status                 SMALLINT              NOT NULL,
    tokens                 BIGINT,
    delegator_shares       DECIMAL               NOT NULL,
    description_id         INT references validator_description (id),
    unbonding_height       BIGINT                NOT NULL,
    unbonding_time         DATE                  NOT NULL,
    commissions            INT                   NOT NULL references validator_commission (id),
    min_self_delegation    INT                   NOT NULL,
    self_delegation_ration DECIMAL               NOT NULL
);

CREATE TABLE delegation
(
    id                SERIAL PRIMARY KEY,
    delegator_address TEXT    NOT NULL,
    validator_address TEXT    NOT NULL references validator_staking_info (operator_address),
    shares            DECIMAL NOT NULL
);

CREATE TABLE delegation_balance
(
    delegation_id INT     NOT NULL references delegation (id),
    denom         TEXT    NOT NULL,
    amount        DECIMAL NOT NULL
)
