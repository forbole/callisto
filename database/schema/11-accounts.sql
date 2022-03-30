-- +migrate Up
CREATE TABLE IF NOT EXISTS vipcoin_chain_accounts_affiliates(
    id               SERIAL PRIMARY KEY NOT NULL,                -- affiliates id
    address          TEXT               NOT NULL,                -- affiliates address
    affiliation_kind INT,                                        -- affiliates AffiliationKind
    extras           JSONB                                       -- affiliates extras
);

CREATE TABLE IF NOT EXISTS vipcoin_chain_accounts_accounts (
    address    TEXT      NOT NULL PRIMARY KEY,                   -- accounts address
    hash       TEXT      NOT NULL,                               -- accounts hash
    public_key TEXT      NOT NULL,                               -- accounts public_key
    kinds      INT[],                                            -- accounts kinds
    state      INT,                                              -- accounts state
    extras     JSONB,                                            -- accounts extras
    affiliates BIGINT[],                                         -- accounts affiliates
    wallets    TEXT[]                                            -- accounts wallets
);

CREATE TABLE IF NOT EXISTS vipcoin_chain_accounts_set_kinds (
    creator  TEXT  NOT NULL,                                     -- set kinds creator
    hash     TEXT  NOT NULL,                                     -- set kinds hash
    kinds    INT[]                                               -- set kinds kinds
);

CREATE TABLE IF NOT EXISTS vipcoin_chain_accounts_set_affiliate_address (
    creator      TEXT  NOT NULL,                                 -- set affiliate address creator
    hash         TEXT  NOT NULL,                                 -- set affiliate address hash
    old_address  TEXT  NOT NULL,                                 -- set affiliate address old_address
    new_address  TEXT  NOT NULL                                  -- set affiliate address new_address
);

CREATE TABLE IF NOT EXISTS vipcoin_chain_accounts_register_user (
    creator                   TEXT  NOT NULL,                    -- register user creator
    address                   TEXT  NOT NULL,                    -- register user address
    hash                      TEXT  NOT NULL,                    -- register user hash
    public_key                TEXT  NOT NULL,                    -- register user public_key
    holder_wallet             TEXT,                              -- register user holder_wallet
    ref_reward_wallet         TEXT,                              -- register user ref_reward_wallet
    holder_wallet_extras      JSONB,                             -- register user holder_wallet_extras
    ref_reward_wallet_extras  JSONB,                             -- register user ref_reward_wallet_extras
    referrer_hash             TEXT                               -- register user referrer_hash
);

CREATE TABLE IF NOT EXISTS vipcoin_chain_accounts_account_migrate (
    creator    TEXT  NOT NULL,                                   -- account migrate creator
    address    TEXT  NOT NULL,                                   -- account migrate address
    hash       TEXT  NOT NULL,                                   -- account migrate hash
    public_key TEXT  NOT NULL,                                   -- accounts public_key
);

CREATE TABLE IF NOT EXISTS vipcoin_chain_accounts_set_affiliate_extra (
    creator           TEXT  NOT NULL,                            -- set affiliate extra creator
    account_hash      TEXT  NOT NULL,                            -- set affiliate extra account_hash
    affiliation_hash  TEXT  NOT NULL,                            -- set affiliate extra affiliation_hash
    extras            JSONB                                      -- set affiliate extra extras
);

CREATE TABLE IF NOT EXISTS vipcoin_chain_accounts_set_extra (
    creator           TEXT  NOT NULL,                            -- set extra creator
    hash              TEXT  NOT NULL,                            -- set extra hash
    extras            JSONB                                      -- set extra extras
);

CREATE TABLE IF NOT EXISTS vipcoin_chain_accounts_set_state (
    creator           TEXT  NOT NULL,                            -- set state creator
    hash              TEXT  NOT NULL,                            -- set state hash
    state             INT   NOT NULL,                            -- set state state
    reason            TEXT  NOT NULL                             -- set state reason
);

CREATE TABLE IF NOT EXISTS vipcoin_chain_accounts_add_affiliate (
    creator           TEXT  NOT NULL,                            -- add affiliate creator
    account_hash      TEXT  NOT NULL,                            -- add affiliate account_hash
    affiliation_hash  TEXT  NOT NULL,                            -- add affiliate affiliation_hash
    affiliation       INT   NOT NULL,                            -- add affiliate affiliation
    extras            JSONB                                      -- add affiliate extras
);

CREATE TABLE IF NOT EXISTS vipcoin_chain_accounts_create_account (
    creator           TEXT  NOT NULL,                            -- create account creator
    hash              TEXT  NOT NULL,                            -- create account hash
    address           TEXT  NOT NULL,                            -- create account address
    public_key        TEXT  NOT NULL,                            -- create account public_key
    kinds             INT[],                                     -- create account kinds
    state             INT,                                       -- create account state
    extras            JSONB                                      -- create account extras
);


-- +migrate Down
DROP TABLE IF EXISTS vipcoin_chain_accounts_affiliates CASCADE;
DROP TABLE IF EXISTS vipcoin_chain_accounts_accounts CASCADE;
DROP TABLE IF EXISTS vipcoin_chain_accounts_set_kinds CASCADE;
DROP TABLE IF EXISTS vipcoin_chain_accounts_set_affiliate_address CASCADE;
DROP TABLE IF EXISTS vipcoin_chain_accounts_register_user CASCADE;
DROP TABLE IF EXISTS vipcoin_chain_accounts_account_migrate CASCADE;
DROP TABLE IF EXISTS vipcoin_chain_accounts_set_affiliate_extra CASCADE;
DROP TABLE IF EXISTS vipcoin_chain_accounts_set_extra CASCADE;
DROP TABLE IF EXISTS vipcoin_chain_accounts_set_state CASCADE;
DROP TABLE IF EXISTS vipcoin_chain_accounts_add_affiliate CASCADE;
DROP TABLE IF EXISTS vipcoin_chain_accounts_create_account CASCADE;
