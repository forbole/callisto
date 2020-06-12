CREATE TABLE staking_pool
(
    height            BIGINT                      NOT NULL,
    timestamp         TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    bonded_tokens     BIGINT                      NOT NULL,
    not_bonded_tokens BIGINT                      NOT NULL
);

CREATE TABLE validator_uptime
(
    height                BIGINT                NOT NULL,
    validator_address     CHARACTER VARYING(52) NOT NULL REFERENCES validator (consensus_address),
    signed_blocks_window  BIGINT                NOT NULL,
    missed_blocks_counter BIGINT                NOT NULL
);

CREATE TABLE validator_info
(
    consensus_address CHARACTER VARYING(52) NOT NULL references validator (consensus_address),
    operator_address  TEXT                  NOT NULL
);

