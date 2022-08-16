CREATE TABLE market_params
(
    one_row_id BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    params     JSONB   NOT NULL,
    height     BIGINT  NOT NULL,
    CHECK (one_row_id)
);

CREATE TYPE LEASE_ID AS
(
    owner       TEXT,
    d_seq       TEXT,
    g_seq       TEXT,
    o_seq       TEXT,
    provider    TEXT
);


CREATE TABLE akash_lease
(
    /* ---- lease id is consist of owner, d_seq, g_seq, o_seq, and provider ---- */
    owner           TEXT        NOT NULL,
    d_seq           TEXT        NOT NULL,
    g_seq           TEXT        NOT NULL,
    o_seq           TEXT        NOT NULL,
    provider        TEXT        NOT NULL,

    lease_state     INT         NOT NULL,
    price           DEC_COIN    NOT NULL DEFAULT'{}'::JSONB

    CONSTRAINT unique_lease_id UNIQUE (owner, d_seq, g_seq, o_seq, provider)
);
CREATE INDEX provider_address_index ON provider (owner_address);
