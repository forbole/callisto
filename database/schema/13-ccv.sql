/* ---- PARAMS ---- */

CREATE TABLE ccv_provider_params
(
    one_row_id BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    params     JSONB   NOT NULL,
    height     BIGINT  NOT NULL,
    CHECK (one_row_id)
);

CREATE TABLE ccv_consumer_params
(
    one_row_id BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    params     JSONB   NOT NULL,
    height     BIGINT  NOT NULL,
    CHECK (one_row_id)
);


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
