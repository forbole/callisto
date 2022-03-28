-- +migrate Up
CREATE TABLE IF NOT EXISTS vipcoin_chain_banking_transfers (
    id                  SERIAL PRIMARY KEY NOT NULL,            -- banking id
    asset               TEXT NOT NULL,                          -- banking asset
    amount              BIGSERIAL,                              -- banking amount
    kind                INT,                                    -- banking kind
    extras              JSONB,                                  -- banking extras
    timestamp           TIMESTAMP WITHOUT TIME ZONE NOT NULL,   -- banking timestamp
    tx_hash             TEXT NOT NULL                           -- banking tx hash
);

CREATE TABLE IF NOT EXISTS vipcoin_chain_banking_system_transfer (
    creator             TEXT NOT NULL,                      -- banking creator
    wallet_from         TEXT NOT NULL,                      -- banking wallet from
    wallet_to           TEXT NOT NULL,                      -- banking wallet to
    asset               TEXT NOT NULL,                      -- banking asset
    amount              BIGSERIAL,                          -- banking amount
    extras              JSONB                               -- banking extras
);

CREATE TABLE IF NOT EXISTS vipcoin_chain_banking_reward_manager_address (
    address            TEXT NOT NULL                        -- banking address
);

CREATE TABLE IF NOT EXISTS vipcoin_chain_banking_set_reward_manager_address (
    creator            TEXT NOT NULL,                       -- banking creator
    address            TEXT NOT NULL                        -- banking address
);

CREATE TABLE IF NOT EXISTS vipcoin_chain_banking_system_reward_transfer (
    creator             TEXT NOT NULL,                      -- banking creator
    wallet_from         TEXT NOT NULL,                      -- banking wallet from
    wallet_to           TEXT NOT NULL,                      -- banking wallet to
    asset               TEXT NOT NULL,                      -- banking asset
    amount              BIGSERIAL,                          -- banking amount
    extras              JSONB                               -- banking extras
);

CREATE TABLE IF NOT EXISTS vipcoin_chain_banking_payment (
    creator             TEXT NOT NULL,                      -- banking creator
    wallet_from         TEXT NOT NULL,                      -- banking wallet from
    wallet_to           TEXT NOT NULL,                      -- banking wallet to
    asset               TEXT NOT NULL,                      -- banking asset
    amount              BIGSERIAL,                          -- banking amount
    extras              JSONB                               -- banking extras
);

CREATE TABLE IF NOT EXISTS vipcoin_chain_banking_withdraw (
    creator             TEXT NOT NULL,                      -- banking creator
    wallet              TEXT NOT NULL,                      -- banking wallet
    asset               TEXT NOT NULL,                      -- banking asset
    amount              BIGSERIAL,                          -- banking amount
    extras              JSONB                               -- banking extras
);

CREATE TABLE IF NOT EXISTS vipcoin_chain_banking_issue (
    creator             TEXT NOT NULL,                      -- banking creator
    wallet              TEXT NOT NULL,                      -- banking wallet
    asset               TEXT NOT NULL,                      -- banking asset
    amount              BIGSERIAL,                          -- banking amount
    extras              JSONB                               -- banking extras
);

-- +migrate Down
DROP TABLE IF EXISTS vipcoin_chain_banking_transfers CASCADE;
DROP TABLE IF EXISTS vipcoin_chain_banking_system_transfer CASCADE;
DROP TABLE IF EXISTS vipcoin_chain_banking_reward_manager_address CASCADE;
DROP TABLE IF EXISTS vipcoin_chain_banking_set_reward_manager_address CASCADE;
DROP TABLE IF EXISTS vipcoin_chain_banking_system_reward_transfer CASCADE;
DROP TABLE IF EXISTS vipcoin_chain_banking_payment CASCADE;
DROP TABLE IF EXISTS vipcoin_chain_banking_withdraw CASCADE;
DROP TABLE IF EXISTS vipcoin_chain_banking_issue CASCADE;
