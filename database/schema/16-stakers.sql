/* ---- PARAMS ---- */
CREATE TABLE stakers_params
(
    one_row_id BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    params     JSONB   NOT NULL,
    height     BIGINT  NOT NULL,
    CHECK (one_row_id)
);

/* ---- PROTOCOL VALIDATOR ---- */

CREATE TABLE protocol_validator
(
    address TEXT   NOT NULL REFERENCES account (address) PRIMARY KEY,
    height  BIGINT NOT NULL
);
CREATE INDEX protocol_validator_height_index ON protocol_validator (height);

CREATE TABLE protocol_validator_commission
(
    address                     TEXT    NOT NULL REFERENCES protocol_validator (address) PRIMARY KEY,
    commission                  TEXT    NOT NULL,
    pending_commission_change   TEXT    NOT NULL DEFAULT '{}'::JSONB,
    self_delegation             BIGINT  NOT NULL,
    height                      BIGINT  NOT NULL
);
CREATE INDEX protocol_validator_commission_height_index ON protocol_validator_commission (height);

CREATE TABLE protocol_validator_delegation
(
    address             TEXT    NOT NULL REFERENCES protocol_validator (address) PRIMARY KEY,
    self_delegation     BIGINT  NOT NULL,
    total_delegation    BIGINT  NOT NULL,
    delegator_count     BIGINT  NOT NULL,
    height              BIGINT  NOT NULL
);
CREATE INDEX protocol_validator_delegation_height_index ON protocol_validator_delegation (height);

CREATE TABLE protocol_validator_description
(
    address           TEXT   NOT NULL REFERENCES protocol_validator (address) PRIMARY KEY,
    moniker           TEXT,
    identity          TEXT,
    avatar_url        TEXT,
    website           TEXT,
    security_contact  TEXT,
    details           TEXT,
    height            BIGINT NOT NULL
);
CREATE INDEX protocol_validator_description_height_index ON protocol_validator_description (height);

CREATE TABLE protocol_validator_pool
(
    id                  SERIAL  NOT NULL PRIMARY KEY,
    address             TEXT    NOT NULL REFERENCES protocol_validator (address),
    validator_address   TEXT    NOT NULL,
    balance             BIGINT  NOT NULL,
    pool                TEXT    NOT NULL REFERENCES pool (name),
    height              BIGINT  NOT NULL,
    CONSTRAINT unique_protocol_validator_pool UNIQUE (address, pool) 

);
CREATE INDEX protocol_validator_pool_height_index ON protocol_validator_pool (height);
