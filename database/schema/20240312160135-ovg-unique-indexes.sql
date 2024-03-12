-- +migrate Up
CREATE UNIQUE INDEX idx_overgold_allowed_addresses ON overgold_allowed_addresses (creator, address);
CREATE UNIQUE INDEX idx_overgold_allowed_create_addresses ON overgold_allowed_create_addresses (tx_hash, creator);
CREATE UNIQUE INDEX idx_overgold_allowed_delete_by_addresses ON overgold_allowed_delete_by_addresses (tx_hash, creator);
CREATE UNIQUE INDEX idx_overgold_allowed_delete_by_id ON overgold_allowed_delete_by_id (tx_hash, creator);

CREATE UNIQUE INDEX idx_overgold_allowed_update_addresses ON overgold_allowed_update_addresses (tx_hash, creator);
CREATE UNIQUE INDEX idx_msg_multi_send ON msg_multi_send (tx_hash);
CREATE UNIQUE INDEX idx_msg_send ON msg_send (tx_hash);

CREATE UNIQUE INDEX idx_overgold_core_issue_tx_hash ON overgold_core_issue (tx_hash);
CREATE UNIQUE INDEX idx_overgold_core_withdraw_tx_hash ON overgold_core_withdraw (tx_hash);
CREATE UNIQUE INDEX idx_overgold_core_send_tx_hash ON overgold_core_send (tx_hash);

CREATE UNIQUE INDEX idx_overgold_referral_set_referrer_tx_hash ON overgold_referral_set_referrer (tx_hash);

CREATE UNIQUE INDEX idx_overgold_stake_sell_tx_hash ON overgold_stake_sell (tx_hash);
CREATE UNIQUE INDEX idx_overgold_stake_buy_tx_hash ON overgold_stake_buy (tx_hash);
CREATE UNIQUE INDEX idx_overgold_stake_sell_cancel_tx_hash ON overgold_stake_sell_cancel (tx_hash);

CREATE UNIQUE INDEX idx_overgold_stake_distribute_rewards_tx_hash ON overgold_stake_distribute_rewards (tx_hash);
CREATE UNIQUE INDEX idx_overgold_stake_claim_reward_tx_hash ON overgold_stake_claim_reward (tx_hash);

CREATE UNIQUE INDEX idx_overgold_feeexcluder_create_address_tx_hash ON overgold_feeexcluder_create_address (tx_hash);
CREATE UNIQUE INDEX idx_overgold_feeexcluder_update_address_tx_hash ON overgold_feeexcluder_update_address (tx_hash);
CREATE UNIQUE INDEX idx_overgold_feeexcluder_delete_address_tx_hash ON overgold_feeexcluder_delete_address (tx_hash);
CREATE UNIQUE INDEX idx_overgold_feeexcluder_create_tariffs_tx_hash ON overgold_feeexcluder_create_tariffs (tx_hash);
CREATE UNIQUE INDEX idx_overgold_feeexcluder_update_tariffs_tx_hash ON overgold_feeexcluder_update_tariffs (tx_hash);
CREATE UNIQUE INDEX idx_overgold_feeexcluder_delete_tariffs_tx_hash ON overgold_feeexcluder_delete_tariffs (tx_hash);
CREATE UNIQUE INDEX idx_overgold_feeexcluder_daily_stats ON overgold_feeexcluder_daily_stats (msg_id);
CREATE UNIQUE INDEX idx_overgold_feeexcluder_fees ON overgold_feeexcluder_fees (msg_id, amount_from, fee);
CREATE UNIQUE INDEX idx_overgold_feeexcluder_genesis_state ON overgold_feeexcluder_genesis_state (address_count, daily_stats_count);
CREATE UNIQUE INDEX idx_overgold_feeexcluder_address ON overgold_feeexcluder_address (address);
CREATE UNIQUE INDEX idx_overgold_feeexcluder_stats ON overgold_feeexcluder_stats (daily_stats_id);
CREATE UNIQUE INDEX idx_overgold_feeexcluder_tariffs ON overgold_feeexcluder_tariffs (denom);
CREATE UNIQUE INDEX idx_overgold_feeexcluder_tariff ON overgold_feeexcluder_tariff (msg_id, amount, denom);
CREATE UNIQUE INDEX idx_overgold_feeexcluder_m2m_genesis_state_address ON overgold_feeexcluder_m2m_genesis_state_address (genesis_state_id, address_id);
CREATE UNIQUE INDEX idx_overgold_feeexcluder_m2m_genesis_state_stats ON overgold_feeexcluder_m2m_genesis_state_stats (genesis_state_id, stats_id);
CREATE UNIQUE INDEX idx_overgold_feeexcluder_m2m_genesis_state_daily_stats ON overgold_feeexcluder_m2m_genesis_state_daily_stats (genesis_state_id, daily_stats_id);
CREATE UNIQUE INDEX idx_overgold_feeexcluder_m2m_genesis_state_tariffs ON overgold_feeexcluder_m2m_genesis_state_tariffs (genesis_state_id, tariffs_id);


CREATE UNIQUE INDEX idx_overgold_stake_transfer_from_user_tx_hash ON overgold_stake_transfer_from_user (tx_hash);
CREATE UNIQUE INDEX idx_overgold_stake_transfer_to_user_tx_hash ON overgold_stake_transfer_to_user (tx_hash);

-- +migrate Down

DROP INDEX IF EXISTS idx_overgold_allowed_addresses;
DROP INDEX IF EXISTS idx_overgold_allowed_create_addresses;
DROP INDEX IF EXISTS idx_overgold_allowed_delete_by_addresses;
DROP INDEX IF EXISTS idx_overgold_allowed_delete_by_id;
DROP INDEX IF EXISTS idx_overgold_allowed_update_addresses;

DROP INDEX IF EXISTS idx_msg_multi_send;
DROP INDEX IF EXISTS idx_msg_send;

DROP INDEX IF EXISTS idx_overgold_core_issue_tx_hash;
DROP INDEX IF EXISTS idx_overgold_core_withdraw_tx_hash;
DROP INDEX IF EXISTS idx_overgold_core_send_tx_hash;

DROP INDEX IF EXISTS idx_overgold_referral_set_referrer_tx_hash;

DROP INDEX IF EXISTS idx_overgold_stake_sell_tx_hash;
DROP INDEX IF EXISTS idx_overgold_stake_buy_tx_hash;
DROP INDEX IF EXISTS idx_overgold_stake_sell_cancel_tx_hash;
DROP INDEX IF EXISTS idx_overgold_stake_distribute_rewards_tx_hash;
DROP INDEX IF EXISTS idx_overgold_stake_claim_reward_tx_hash;

DROP INDEX IF EXISTS idx_overgold_feeexcluder_create_address_tx_hash;
DROP INDEX IF EXISTS idx_overgold_feeexcluder_update_address_tx_hash;
DROP INDEX IF EXISTS idx_overgold_feeexcluder_delete_address_tx_hash;
DROP INDEX IF EXISTS idx_overgold_feeexcluder_create_tariffs_tx_hash;
DROP INDEX IF EXISTS idx_overgold_feeexcluder_update_tariffs_tx_hash;
DROP INDEX IF EXISTS idx_overgold_feeexcluder_delete_tariffs_tx_hash;
DROP INDEX IF EXISTS idx_overgold_feeexcluder_daily_stats;
DROP INDEX IF EXISTS idx_overgold_feeexcluder_fees;
DROP INDEX IF EXISTS idx_overgold_feeexcluder_genesis_state;
DROP INDEX IF EXISTS idx_overgold_feeexcluder_address;
DROP INDEX IF EXISTS idx_overgold_feeexcluder_stats;
DROP INDEX IF EXISTS idx_overgold_feeexcluder_tariffs;
DROP INDEX IF EXISTS idx_overgold_feeexcluder_tariff;
DROP INDEX IF EXISTS idx_overgold_feeexcluder_m2m_genesis_state_address;
DROP INDEX IF EXISTS idx_overgold_feeexcluder_m2m_genesis_state_stats;
DROP INDEX IF EXISTS idx_overgold_feeexcluder_m2m_genesis_state_daily_stats;
DROP INDEX IF EXISTS idx_overgold_feeexcluder_m2m_genesis_state_tariffs;


DROP INDEX IF EXISTS idx_overgold_stake_transfer_from_user_tx_hash;
DROP INDEX IF EXISTS idx_overgold_stake_transfer_to_user_tx_hash;
