/* ---- PARAMS ---- */
CREATE TABLE ccv_consumer_params
(
    one_row_id BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    params     JSONB   NOT NULL,
    height     BIGINT  NOT NULL,
    CHECK (one_row_id)
);

CREATE TABLE ccv_provider_params
(
    one_row_id BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    params     JSONB   NOT NULL,
    height     BIGINT  NOT NULL,
    CHECK (one_row_id)
);

/* ---- CCV CONSUMER CHAIN STATE ---- */
CREATE TABLE ccv_consumer_chain
(
    provider_client_id             TEXT UNIQUE NOT NULL,
    provider_channel_id            TEXT,
    new_chain                      BOOLEAN NOT NULL DEFAULT TRUE,
    provider_client_state          JSONB,
    provider_consensus_state       JSONB,
    maturing_packets               JSONB,
    initial_val_set                JSONB,
    height_to_valset_update_id     JSONB,
    outstanding_downtime_slashing  JSONB,
    pending_consumer_packets       JSONB,
    last_transmission_block_height JSONB,
    height                         BIGINT  NOT NULL
);
CREATE INDEX ccv_consumer_chain_height_index ON ccv_consumer_chain (height);
CREATE INDEX ccv_consumer_chain_provider_client_id_index ON ccv_consumer_chain (provider_client_id);


/* ---- CCV PROVIDER CHAIN STATE ---- */
CREATE TABLE ccv_provider_chain
(
    valset_update_id             INTEGER,
    consumer_states              JSONB,
    unbonding_ops                JSONB,
    mature_unbonding_ops         JSONB,
    valset_update_id_to_height   JSONB,
    consumer_addition_proposals  JSONB,
    consumer_removal_proposals   JSONB,
    validator_consumer_pubkeys   JSONB,
    validators_by_consumer_addr  JSONB,
    consumer_addrs_to_prune      JSONB,
    height                       BIGINT  NOT NULL
);
CREATE INDEX ccv_provider_chain_height_index ON ccv_provider_chain (height);


/* ---- CCV FEE DISTRIBUTION ---- */
CREATE TABLE ccv_fee_distribution
(
    current_height         BIGINT NOT NULL,
    last_height            BIGINT,
    next_height            BIGINT NOT NULL,
    distribution_fraction  TEXT,
    total                  TEXT,
    to_provider            TEXT REFERENCES account (address),
    to_consumer            TEXT REFERENCES account (address),
    height                 BIGINT  NOT NULL,
    CONSTRAINT unique_provider_consumer_fee_distribution UNIQUE(to_provider, to_consumer) 

);
CREATE INDEX ccv_fee_distribution_height_index ON ccv_provider_chain (height);

/* ---- CCV PROPOSAL ---- */
CREATE TABLE ccv_proposal
(
    id                                    INTEGER    NOT NULL PRIMARY KEY,
    title                                 TEXT       NOT NULL,
    description                           TEXT       NOT NULL,
    chain_id                              TEXT       NOT NULL,
    genesis_hash                          TEXT       NOT NULL,
    binary_hash                           TEXT       NOT NULL,
    proposal_type                         TEXT       NOT NULL,
    proposal_route                        TEXT       NOT NULL,
    spawn_time                            TIMESTAMP,
    stop_time                             TIMESTAMP,
    initial_height                        BIGINT,
    unbonding_period                      TIMESTAMP  NOT NULL,
    ccv_timeout_period                    TIMESTAMP  NOT NULL,
    consumer_redistribution_fraction      TEXT       NOT NULL, 
    blocks_per_distribution_transmission  BIGINT     NOT NULL,
    historical_entries                    INTEGER    NOT NULL,
    status                                TEXT       NOT NULL,
    submit_time                           TEXT       NOT NULL,
    proposer_address                      TEXT       NOT NULL REFERENCES account (address),
    height                                BIGINT     NOT NULL
);
CREATE INDEX ccv_proposal_proposer_address_index ON ccv_proposal (proposer_address);
