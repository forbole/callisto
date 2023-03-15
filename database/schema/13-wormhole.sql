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

/* ---- GUARDIAN SET  ---- */

CREATE TABLE guardian_set
(
    index              INTEGER NOT NULL,
    keys               JSONB   NOT NULL,
    expiration_time    BIGINT  NOT NULL,
    height             BIGINT  NOT NULL
);
CREATE INDEX guardian_set_keys_index ON guardian_set (keys);

/* ---- ALLOW LIST  ---- */

CREATE TABLE validator_allow_list
(
    validator_address  TEXT    NOT NULL,
    allowed_address    TEXT    NOT NULL,
    name               TEXT    NOT NULL,
    height             BIGINT  NOT NULL
);
CREATE INDEX validator_allow_list_validator_address_index ON validator_allow_list (validator_address);
CREATE INDEX validator_allow_listt_allowed_address_index ON validator_allow_list (allowed_address);
CREATE INDEX validator_allow_list_name_index ON validator_allow_list (name);


