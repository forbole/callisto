-- +migrate Up
CREATE TABLE IF NOT EXISTS vipcoin_chain_assets_assets(
    name           TEXT NOT NULL PRIMARY KEY, -- assets name
    issuer         TEXT NOT NULL,             -- assets issuer
    policies       INT[],                     -- assets policies
    state          INT,                       -- assets state
    issued         BIGINT,                    -- assets issued
    burned         BIGINT,                    -- assets burned
    withdrawn      BIGINT,                    -- assets withdrawn
    in_circulation BIGINT,                    -- assets in_circulation
    precision      INT,                       -- assets precision
    fee_percent    INT,                       -- assets fee_percent
    extras         JSONB                      -- assets extras
);

CREATE TABLE IF NOT EXISTS vipcoin_chain_assets_set_extra(
    id                SERIAL UNIQUE PRIMARY KEY NOT NULL,
    transaction_hash  TEXT NOT NULL REFERENCES transaction (hash),
    creator           TEXT NOT NULL, -- set extra creator
    name              TEXT NOT NULL, -- set extra name
    extras            JSONB          -- set extra extras
);

CREATE TABLE IF NOT EXISTS vipcoin_chain_assets_manage(
    id               SERIAL UNIQUE PRIMARY KEY NOT NULL,
    transaction_hash TEXT NOT NULL REFERENCES transaction (hash),
    creator          TEXT NOT NULL, -- manage assets creator
    name             TEXT NOT NULL, -- manage assets extra name
    policies         INT[],         -- manage assets policies
    state            INT,           -- manage assets state
    precision        INT,           -- assets precision
    fee_percent      INT,           -- assets fee_percent
    issued           BIGINT,        -- manage assets issued
    burned           BIGINT,        -- manage assets burned
    withdrawn        BIGINT,        -- manage assets withdrawn
    in_circulation   BIGINT         -- manage assets in_circulation
);

CREATE TABLE IF NOT EXISTS vipcoin_chain_assets_create(
    id                SERIAL UNIQUE PRIMARY KEY NOT NULL,
    transaction_hash  TEXT NOT NULL REFERENCES transaction (hash),
    creator           TEXT NOT NULL, -- create assets creator
    name              TEXT NOT NULL, -- create assets extra name
    issuer            TEXT NOT NULL, -- create assets issuer
    policies          INT[],         -- create assets policies
    state             INT,           -- create assets state
    precision         INT,           -- assets precision
    fee_percent       INT,           -- assets fee_percent
    extras            JSONB          -- create assets extras
);

-- +migrate Down
DROP TABLE IF EXISTS vipcoin_chain_assets_assets CASCADE;
DROP TABLE IF EXISTS vipcoin_chain_assets_set_extra CASCADE;
DROP TABLE IF EXISTS vipcoin_chain_assets_manage CASCADE;
DROP TABLE IF EXISTS vipcoin_chain_assets_create CASCADE;
