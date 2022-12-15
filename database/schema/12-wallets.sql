-- +migrate Up
CREATE TABLE IF NOT EXISTS overgold_chain_wallets_wallets (
    address             TEXT    NOT NULL PRIMARY KEY,       -- wallet address
    account_address     TEXT    NOT NULL,                   -- account address
    kind                INT,                                -- kind for different owners (issuer, holder, market etc.)
    state               INT,                                -- activity state (unspecified, active, blocked etc.)
    balance             JSONB,                              -- wallet balance
    extras              JSONB,                              -- extras for additional data
    default_status      BOOLEAN,                            -- wallet status for default use
    CONSTRAINT fk_accounts_accounts
      FOREIGN KEY(account_address) 
      REFERENCES overgold_chain_accounts_accounts(address) ON UPDATE CASCADE
);

CREATE TABLE IF NOT EXISTS overgold_chain_wallets_set_wallet_kind (
    id                SERIAL UNIQUE PRIMARY KEY NOT NULL,
    transaction_hash  TEXT    NOT NULL REFERENCES transaction (hash),
    creator           TEXT    NOT NULL,                           -- message creator
    address           TEXT    NOT NULL,                           -- target wallet address
    kind              INT                                         -- new kind for target wallet
);

CREATE TABLE IF NOT EXISTS overgold_chain_wallets_set_wallet_state (
    id               SERIAL UNIQUE PRIMARY KEY NOT NULL,
    transaction_hash TEXT    NOT NULL REFERENCES transaction (hash),
    creator          TEXT    NOT NULL,                           -- message creator
    address          TEXT    NOT NULL,                           -- target wallet address
    state            INT                                         -- new state for target wallet
);

CREATE TABLE IF NOT EXISTS overgold_chain_wallets_create_wallet (
    id                  SERIAL UNIQUE PRIMARY KEY NOT NULL,
    transaction_hash    TEXT    NOT NULL REFERENCES transaction (hash),
    creator             TEXT    NOT NULL,                   -- message creator
    address             TEXT    NOT NULL,                   -- target wallet address
    account_address     TEXT    NOT NULL,                   -- new account address for target wallet
    kind                INT,                                -- new kind for target wallet
    state               INT,                                -- new state for target wallet
    extras              JSONB                               -- new extras for target wallet
);

CREATE TABLE IF NOT EXISTS overgold_chain_wallets_create_wallet_with_balance (
    id                  SERIAL UNIQUE PRIMARY KEY NOT NULL,
    transaction_hash    TEXT        NOT NULL REFERENCES transaction (hash),
    creator             TEXT        NOT NULL,               -- message creator
    address             TEXT        NOT NULL,               -- target wallet address
    account_address     TEXT        NOT NULL,               -- new account address for target wallet
    kind                INT,                                -- new kind for target wallet
    state               INT,                                -- new state for target wallet
    extras              JSONB,                              -- new state for target wallet
    default_status      BOOLEAN,                            -- new default status for target wallet
    balance             JSONB                               -- new balance for target wallet
);

CREATE TABLE IF NOT EXISTS overgold_chain_wallets_set_default_wallet (
    id               SERIAL UNIQUE PRIMARY KEY NOT NULL,
    transaction_hash TEXT        NOT NULL REFERENCES transaction (hash),
    creator          TEXT        NOT NULL,                       -- message creator
    address          TEXT        NOT NULL                        -- target wallet address
);

CREATE TABLE IF NOT EXISTS overgold_chain_wallets_set_extra (
    id                SERIAL UNIQUE PRIMARY KEY NOT NULL,
    transaction_hash  TEXT        NOT NULL REFERENCES transaction (hash),
    creator           TEXT        NOT NULL,                       -- message creator
    address           TEXT        NOT NULL,                       -- target wallet address
    extras            JSONB                                       -- new extras for target wallet
);

-- +migrate Down
DROP TABLE IF EXISTS overgold_chain_wallets_wallets CASCADE;
DROP TABLE IF EXISTS overgold_chain_wallets_set_wallet_kind CASCADE;
DROP TABLE IF EXISTS overgold_chain_wallets_set_wallet_state CASCADE;
DROP TABLE IF EXISTS overgold_chain_wallets_create_wallet CASCADE;
DROP TABLE IF EXISTS overgold_chain_wallets_create_wallet_with_balance CASCADE;
DROP TABLE IF EXISTS overgold_chain_wallets_set_default_wallet CASCADE;
DROP TABLE IF EXISTS overgold_chain_wallets_set_extra CASCADE;
