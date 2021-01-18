CREATE TABLE validator_signing_info
(
    validator_address     TEXT                        NOT NULL UNIQUE REFERENCES validator (consensus_address),
    start_height          BIGINT                      NOT NULL,
    index_offset          BIGINT                      NOT NULL,
    jailed_until          TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    tombstoned            BOOLEAN                     NOT NULL,
    missed_blocks_counter BIGINT                      NOT NULL,
    height                BIGINT                      NOT NULL,
);

CREATE TABLE slashing_params
(
    signed_block_window        BIGINT  NOT NULL,
    min_signed_per_window      BIGINT  NOT NULL,
    downtime_jail_duration     TIMESTAMP WITHOUT TIME ZONE,
    slash_fraction_double_sign DECIMAL NOT NULL,
    slash_fraction_downtime    DECIMAL NOT NULL,
    height                     BIGINT  NOT NULL
);