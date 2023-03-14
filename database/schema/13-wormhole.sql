/* ---- CONFIG ---- */

CREATE TABLE wormhole_config
(
    one_row_id BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    config     JSONB   NOT NULL,
    height     BIGINT  NOT NULL,
    CHECK (one_row_id)
);
CREATE INDEX wormhole_config_height_index ON wormhole_config (height);

/* ---- GUARDIAN VALIDATOR  ---- */

CREATE TABLE guardian_validator
(
    guardian_key       TEXT   NOT NULL UNIQUE PRIMARY KEY,
    validator_address  TEXT   NOT NULL UNIQUE,
    height             BIGINT NOT NULL
);
CREATE INDEX guardian_validator_guardian_key_index ON guardian_validator (guardian_key);
CREATE INDEX guardian_validator_validator_address_index ON guardian_validator (validator_address);
