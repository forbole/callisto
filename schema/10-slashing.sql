CREATE TABLE validator_signing_info
(
    validator_address     TEXT                        NOT NULL REFERENCES validator (consensus_address),
    start_height          BIGINT                      NOT NULL,
    index_offset          BIGINT                      NOT NULL,
    jailed_until          TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    tombstoned            TEXT                        NOT NULL,
    missed_blocks_counter BIGINT                      NOT NULL,
    height                BIGINT                      NOT NULL,
    UNIQUE (validator_address, height)
);