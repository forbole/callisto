CREATE TYPE ACCOUNT_ID AS
(
    scope   TEXT,
    xid     TEXT
);


CREATE TABLE lease
(
    /* Lease ID */
    owner_address       TEXT        NOT NULL REFERENCES account (address),
    dseq                BIGINT      NOT NULL,
    gseq                BIGINT      NOT NULL,
    oseq                BIGINT      NOT NULL,
    provider_address    TEXT        NOT NULL REFERENCES account (address),
    CONSTRAINT unique_lease_id UNIQUE(owner_address, dseq, gseq, oseq, provider_address),
    
    /* Lease */
    lease_state         TEXT        NOT NULL,
    price               DEC_COIN    NOT NULL,
    created_at          BIGINT      NOT NULL,
    closed_on           BIGINT      NOT NULL,

    /* Escrow Payment */
    account_id      ACCOUNT_ID  NOT NULL,
    payment_id      TEXT        NOT NULL,
    payment_state   INT         NOT NULL,
    rate            DEC_COIN    NOT NULL,
    balance         DEC_COIN    NOT NULL,
    withdrawn       COIN        NOT NULL,

    height          BIGINT      NOT NULL
);
CREATE INDEX lease_height_index ON lease (height);
