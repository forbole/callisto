CREATE TABLE validator_signing_info
(
    validator_address     TEXT                        NOT NULL,
    start_height          BIGINT                      NOT NULL,
    index_offset          BIGINT                      NOT NULL,
    jailed_until          TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    tombstoned            BOOLEAN                     NOT NULL,
    missed_blocks_counter BIGINT                      NOT NULL,
    height                BIGINT                      NOT NULL REFERENCES block (height),
    UNIQUE (validator_address, height)
);
CREATE INDEX validator_signing_info_validator_address_index ON validator_signing_info (validator_address);
CREATE INDEX validator_signing_info_height_index ON validator_signing_info (height);

CREATE TABLE slashing_params
(
    signed_block_window        BIGINT  NOT NULL,
    min_signed_per_window      DECIMAL NOT NULL,
    downtime_jail_duration     BIGINT  NOT NULL,
    slash_fraction_double_sign DECIMAL NOT NULL,
    slash_fraction_downtime    DECIMAL NOT NULL,
    height                     BIGINT  NOT NULL REFERENCES block (height)
);
CREATE INDEX slashing_params_height_index ON slashing_params (height);