-- +migrate Up
CREATE TABLE IF NOT EXISTS vipcoin_chain_wallets_wallets (
    address             TEXT    NOT NULL PRIMARY KEY,       -- wallets address
    account_address     TEXT    NOT NULL,                   -- wallets account address
    kind                INT,                                -- wallets kind
    state               INT,                                -- wallets state
    balance             JSONB,                              -- wallets balance
    extras              JSONB,                              -- wallets extras
    default_status      BOOLEAN                             -- wallets default
);

CREATE TABLE IF NOT EXISTS vipcoin_chain_wallets_set_wallet_kind (
    creator     TEXT    NOT NULL,                           -- set wallet kind creator
    address     TEXT    NOT NULL,                           -- set wallet kind address
    kind        INT                                         -- set wallet kind kind
);

CREATE TABLE IF NOT EXISTS vipcoin_chain_wallets_set_wallet_state (
    creator     TEXT    NOT NULL,                           -- set wallet state creator
    address     TEXT    NOT NULL,                           -- set wallet state address
    state       INT                                         -- set wallet state state
);

CREATE TABLE IF NOT EXISTS vipcoin_chain_wallets_create_wallet (
    creator             TEXT    NOT NULL,                   -- create wallet creator
    address             TEXT    NOT NULL,                   -- create wallet address
    account_address     TEXT    NOT NULL,                   -- create wallet account address
    kind                INT,                                -- create wallet kind
    state               INT,                                -- create wallet state
    extras              JSONB                               -- create wallet extras
);

CREATE TABLE IF NOT EXISTS vipcoin_chain_wallets_create_wallet_with_balance (
    creator             TEXT        NOT NULL,               -- create wallet with balance creator
    address             TEXT        NOT NULL,               -- create wallet with balance address
    account_address     TEXT        NOT NULL,               -- create wallet with balance account address
    kind                INT,                                -- create wallet with balance kind
    state               INT,                                -- create wallet with balance state
    extras              JSONB,                              -- create wallet with balance extras
    default_status      BOOLEAN,                            -- create wallet with balance default
    balance             JSONB                               -- create wallet with balance balance
);

CREATE TABLE IF NOT EXISTS vipcoin_chain_wallets_set_default_wallet (
    creator     TEXT        NOT NULL,                       -- set default wallet creator
    address     TEXT        NOT NULL                        -- set default wallet address
);

CREATE TABLE IF NOT EXISTS vipcoin_chain_wallets_set_extra (
    creator     TEXT        NOT NULL,                       -- set extra creator
    address     TEXT        NOT NULL,                       -- set extra address
    extras      JSONB                                       -- set extra extras
);

-- +migrate Down
DROP TABLE IF EXISTS vipcoin_chain_wallets_wallets CASCADE;
DROP TABLE IF EXISTS vipcoin_chain_wallets_set_wallet_kind CASCADE;
DROP TABLE IF EXISTS vipcoin_chain_wallets_set_wallet_state CASCADE;
DROP TABLE IF EXISTS vipcoin_chain_wallets_create_create_wallet_with_balance CASCADE;
DROP TABLE IF EXISTS vipcoin_chain_wallets_create_wallet_with_balance CASCADE;
DROP TABLE IF EXISTS vipcoin_chain_wallets_set_default_wallet CASCADE;
DROP TABLE IF EXISTS vipcoin_chain_wallets_set_extra CASCADE;
