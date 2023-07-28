/* ---- CCV VALIDATOR ---- */
CREATE TABLE ccv_validator
(
    consumer_consensus_address      TEXT   NOT NULL PRIMARY KEY, /* Validator consensus address on consumer chain */
    consumer_self_delegate_address  TEXT   NOT NULL, /* Validator self delegate address on consumer chain */
    consumer_operator_address       TEXT   NOT NULL, /* Validator operator address on consumer chain */
    provider_consensus_address      TEXT   NOT NULL, /* Validator consensus address on provider chain */
    provider_self_delegate_address  TEXT   NOT NULL, /* Validator self delegate address on provider chain */
    provider_operator_address       TEXT   NOT NULL, /* Validator operator address on provider chain */
    height                          BIGINT NOT NULL
);

/* ---- PARAMS ---- */
CREATE TABLE ccv_consumer_params
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
CREATE INDEX ccv_consumer_chain_provider_channel_id_index ON ccv_consumer_chain (provider_channel_id);
CREATE INDEX ccv_consumer_chain_id_index ON ccv_consumer_chain (chain_id);
