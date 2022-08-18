CREATE TABLE market_params
(
    one_row_id BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    params     JSONB   NOT NULL,
    height     BIGINT  NOT NULL,
    CHECK (one_row_id)
);

CREATE TABLE akash_lease
(
    /* ---- lease id is composed of owner, d_seq, g_seq, o_seq, and provider ---- */
    owner           TEXT        NOT NULL,
    d_seq           TEXT        NOT NULL,
    g_seq           TEXT        NOT NULL,
    o_seq           TEXT        NOT NULL,
    provider        TEXT        NOT NULL,

    lease_state     INT         NOT NULL,
    price           DEC_COIN[]  NOT NULL,
    created_at      BIGINT      NOT NULL,
    closed_on       BIGINT      NOT NULL,
    height          BIGINT      NOT NULL,

    CONSTRAINT unique_lease_id UNIQUE (owner, d_seq, g_seq, o_seq, provider)
);
