CREATE TABLE lease_id
(
    id                  SERIAL      NOT NULL PRIMARY KEY,
    owner_address       TEXT        NOT NULL REFERENCES account (address),
    dseq                BIGINT      NOT NULL,
    gseq                BIGINT      NOT NULL,
    oseq                BIGINT      NOT NULL,
    provider_address    TEXT        NOT NULL REFERENCES account (address),
    UNIQUE(owner_address, dseq, gseq, oseq, provider_address)
);

CREATE TABLE lease
(
    lease_id            BIGINT      NOT NULL PRIMARY KEY REFERENCES lease_id (id),
    lease_state         TEXT        NOT NULL,
    price               DEC_COIN    NOT NULL,
    created_at          BIGINT      NOT NULL,
    closed_on           BIGINT      NOT NULL,
    height              BIGINT      NOT NULL
);
CREATE INDEX lease_height_index ON lease (height);


CREATE TYPE ACCOUNT_ID AS
(
    scope   TEXT,
    xid     TEXT
);

CREATE TABLE escrow_payment
(
    lease_id        BIGINT      NOT NULL PRIMARY KEY REFERENCES lease_id (id),
    account_id      ACCOUNT_ID  NOT NULL,
    payment_id      TEXT        NOT NULL,
    owner_address   TEXT        NOT NULL,
    payment_state   INT         NOT NULL,
    rate            DEC_COIN    NOT NULL,
    balance         DEC_COIN    NOT NULL,
    withdrawn       COIN        NOT NULL,
    height          BIGINT      NOT NULL
);

CREATE TABLE market_params
(
    one_row_id          BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    bid_min_deposit     COIN    NOT NULL,
    order_max_bids      INT     NOT NULL,
    CHECK (one_row_id)
);
CREATE INDEX market_params_height_index ON market_params (height);

