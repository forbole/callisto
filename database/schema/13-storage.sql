/* ---- PARAMS ---- */

CREATE TABLE storage_params
(
    one_row_id BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    params     JSONB   NOT NULL,
    height     BIGINT  NOT NULL,
    CHECK (one_row_id)
);

CREATE TABLE storage_providers
(
    address            TEXT     NOT NULL REFERENCES account (address) PRIMARY KEY,
    ip                 TEXT     NOT NULL, 
    total_space        TEXT     NOT NULL,
    burned_contracts   TEXT     NOT NULL,
    creator            TEXT     NOT NULL,
    keybase_identity   TEXT,
    auth_claimers      TEXT[],
    height             BIGINT   NOT NULL
);
CREATE INDEX storage_providers_address_index ON storage_providers (address);
CREATE INDEX creator_index ON storage_providers (creator);
