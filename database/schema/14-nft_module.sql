
CREATE TABLE nft_denom
(
    transaction_hash TEXT NOT NULL REFERENCES transaction (hash),
    id TEXT NOT NULL,
    name TEXT NOT NULL,
    schema TEXT NOT NULL,
    sender TEXT NOT NULL,
    url TEXT NOT NULL,
    PRIMARY KEY(id)
);

CREATE INDEX nft_denom_sender_index ON nft_denom (sender);

CREATE TABLE nft_nft
(
    transaction_hash TEXT NOT NULL REFERENCES transaction (hash),
    id TEXT NOT NULL,
    denom_id TEXT NOT NULL REFERENCES nft_denom (id),
    name TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    uri TEXT NOT NULL,
    tags TEXT[] NOT NULL,
    sender TEXT NOT NULL,
    recipient TEXT NOT NULL,
    uniq_id TEXT NOT NULL,
    PRIMARY KEY(id, denom_id)
);

CREATE INDEX nft_nft_recipient_index ON nft_nft (recipient);

-- TODO: need create FUNCTION for seaching by name/description/tags
