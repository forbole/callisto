CREATE TYPE ACCESS_CONFIG AS
(
    permission  INT,
    address     TEXT
);

CREATE TABLE wasm_code
(
    sender                  TEXT            NOT NULL REFERENCES account (address),
    byte_code               TEXT            NOT NULL,
    instantiate_permission  ACCESS_CONFIG   NULL,
    code_id                 BIGINT          NOT NULL UNIQUE,
    height                  BIGINT          NOT NULL REFERENCES block (height)
);
CREATE INDEX wasm_code_height_index ON wasm_code (height);

CREATE TABLE wasm_contract
(
    sender                  TEXT            NOT NULL REFERENCES account (address),
    admin                   TEXT            NULL REFERENCES account (address),
    code_id                 BIGINT          NOT NULL REFERENCES wasm_code (code_id),
    label                   TEXT            NULL,
    raw_contract_message    JSONB           NOT NULL DEFAULT '{}'::JSONB,
    funds                   COIN[]          NOT NULL DEFAULT '{}',
    contract_address        TEXT            NOT NULL UNIQUE,
    data                    JSONB           NULL,
    instantiated_at         TIMESTAMP       NOT NULL,
    contract_info_extension JSONB           NULL,
    height                  BIGINT          NOT NULL REFERENCES block (height)
);
CREATE INDEX wasm_contract_height_index ON wasm_contract (height);

