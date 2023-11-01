ALTER TABLE nft_nft ADD uniq_id TEXT DEFAULT '';

CREATE INDEX nft_nft_uniq_id_index ON nft_nft (uniq_id);

DELETE FROM workers_storage WHERE key = 'migrate_nfts_worker_nft_migrate_current_height';
DELETE FROM workers_storage WHERE key = 'migrate_nfts_worker_nft_migrate_until_height';
