UPDATE nft_nft
SET uniq_id = concat(id, '@', denom_id);
