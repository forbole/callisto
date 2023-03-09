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
    chain_id                       TEXT,
    provider_client_state          JSONB,
    provider_consensus_state       JSONB,
    initial_val_set                JSONB,
    height                         BIGINT  NOT NULL,
    CONSTRAINT unique_provider_id UNIQUE(provider_client_id, provider_channel_id) 
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

