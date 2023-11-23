CREATE TYPE ACCESS_CONFIG AS
(
    permission  INT,
    address     TEXT
);

CREATE TABLE wasm_params
(
    one_row_id                      BOOLEAN         NOT NULL DEFAULT TRUE PRIMARY KEY,
    code_upload_access              ACCESS_CONFIG   NOT NULL,
    instantiate_default_permission  INT             NOT NULL,
    height                          BIGINT          NOT NULL
);