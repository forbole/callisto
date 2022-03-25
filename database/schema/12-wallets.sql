-- +migrate Up
CREATE TABLE IF NOT EXISTS vipcoin_chain_wallets_wallets (
    address             TEXT    NOT NULL PRIMARY KEY,        -- wallets address
    account_address     TEXT    NOT NULL,                    -- wallets account address
    kind                INT,                                 -- wallets kind
    state               INT,                                 -- wallets state
    balance             JSONB,                               -- wallets balance
    extras              JSONB,                               -- wallets extras
    default_status      BOOLEAN                              -- wallets default
);

CREATE TABLE IF NOT EXISTS vipcoin_chain_wallets_set_wallet_kind (
    creator     TEXT    NOT NULL,                            -- set wallets kind creator
    address     TEXT    NOT NULL,                            -- set wallets kind address
    kind        INT                                          -- set wallets kind kind
);

CREATE TABLE IF NOT EXISTS vipcoin_chain_wallets_set_wallet_state (
    creator     TEXT    NOT NULL,                             -- set wallets state creator
    address     TEXT    NOT NULL,                             -- set wallets state address
    state       INT                                           -- set wallets state state
);

CREATE TABLE IF NOT EXISTS vipcoin_chain_wallets_create_wallet (
    creator             TEXT    NOT NULL,                     -- create wallets creator
    address             TEXT    NOT NULL,                     -- create wallets address
    account_address     TEXT    NOT NULL,                     -- create wallets account address
    kind                INT,                                  -- create wallets kind
    state               INT,                                  -- create wallets state
    extras              JSONB                                 -- create wallets extras
);

CREATE TABLE IF NOT EXISTS vipcoin_chain_wallets_create_wallet_with_balance (
    creator             TEXT        NOT NULL,                 -- create wallets with balance creator
    address             TEXT        NOT NULL,                 -- create wallets with balance address
    account_address     TEXT        NOT NULL,                 -- create wallets with balance account address
    kind                INT,                                  -- create wallets with balance kind
    state               INT,                                  -- create wallets with balance state
    extras              JSONB,                                -- create wallets with balance extras
    default_status      BOOLEAN,                              -- create wallets with balance default
    balance             JSONB                                 -- create wallets with balance balance
);

CREATE TABLE IF NOT EXISTS vipcoin_chain_wallets_set_default_wallet (
    creator     TEXT        NOT NULL,                         -- set default wallets creator
    address     TEXT        NOT NULL                          -- set default wallets address
);

CREATE TABLE IF NOT EXISTS vipcoin_chain_wallets_set_extra (
   creator     TEXT        NOT NULL,                          -- set extra creator
   address     TEXT        NOT NULL,                          -- set extra address
   extras      JSONB                                          -- set extra extras
);

-- +migrate Down
DROP TABLE IF EXISTS vipcoin_chain_wallets_wallets CASCADE;
DROP TABLE IF EXISTS vipcoin_chain_wallets_set_wallet_kind CASCADE;
DROP TABLE IF EXISTS vipcoin_chain_wallets_set_wallet_state CASCADE;
DROP TABLE IF EXISTS vipcoin_chain_wallets_create_create_wallet_with_balance CASCADE;
DROP TABLE IF EXISTS vipcoin_chain_wallets_set_default_wallet CASCADE;
DROP TABLE IF EXISTS vipcoin_chain_wallets_set_extra CASCADE;
