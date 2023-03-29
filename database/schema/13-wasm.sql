CREATE TABLE wasm_params
(
    one_row_id BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    params     JSONB   NOT NULL,
    height     BIGINT  NOT NULL,
    CHECK (one_row_id)
);
CREATE INDEX wasm_params_height_index ON staking_params (height);

CREATE TABLE wasm_code
(
    sender                  TEXT            NULL,
    byte_code               BYTEA           NOT NULL,
    instantiate_permission  JSONB           NULL,
    code_id                 BIGINT          NOT NULL UNIQUE,
    height                  BIGINT          NOT NULL
);
CREATE INDEX wasm_code_height_index ON wasm_code (height);

CREATE TABLE wasm_contract
(
    sender                  TEXT            NULL,
    creator                 TEXT            NOT NULL REFERENCES account (address),
    admin                   TEXT            NULL,
    code_id                 BIGINT          NOT NULL REFERENCES wasm_code (code_id),
    label                   TEXT            NULL,
    raw_contract_message    JSONB           NOT NULL DEFAULT '{}'::JSONB,
    funds                   COIN[]          NOT NULL DEFAULT '{}',
    contract_address        TEXT            NOT NULL UNIQUE,
    data                    TEXT            NULL,
    instantiated_at         TIMESTAMP       NOT NULL,
    contract_info_extension TEXT            NULL,
    contract_states         JSONB           NOT NULL DEFAULT '{}'::JSONB,
    height                  BIGINT          NOT NULL
);
CREATE INDEX wasm_contract_height_index ON wasm_contract (height);

CREATE TABLE wasm_execute_contract
(
    sender                  TEXT            NOT NULL,
    contract_address        TEXT            NOT NULL REFERENCES wasm_contract (contract_address),
    raw_contract_message    JSONB           NOT NULL DEFAULT '{}'::JSONB,
    funds                   COIN[]          NOT NULL DEFAULT '{}',
    data                    TEXT            NULL,
    executed_at             TIMESTAMP       NOT NULL,
    height                  BIGINT          NOT NULL
);
CREATE INDEX execute_contract_height_index ON wasm_execute_contract (height);
 