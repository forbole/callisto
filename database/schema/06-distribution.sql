-- +migrate Up
CREATE TYPE DEC_COIN AS
(
    denom  TEXT,
    amount TEXT
);

/* ---- PARAMS ---- */

CREATE TABLE distribution_params
(
    one_row_id BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    params     JSONB   NOT NULL,
    height     BIGINT  NOT NULL,
    CHECK (one_row_id)
);
CREATE INDEX distribution_params_height_index ON distribution_params (height);


/* ---- COMMUNITY POOL ---- */

CREATE TABLE community_pool
(
    one_row_id bool PRIMARY KEY DEFAULT TRUE,
    coins      DEC_COIN[] NOT NULL,
    height     BIGINT     NOT NULL,
    CONSTRAINT one_row_uni CHECK (one_row_id)
);
CREATE INDEX community_pool_height_index ON community_pool (height);

-- +migrate Down
DROP INDEX IF EXISTS community_pool_height_index;
DROP TABLE IF EXISTS community_pool CASCADE;
DROP INDEX IF EXISTS distribution_params_height_index;
DROP TABLE IF EXISTS distribution_params CASCADE;
DROP TYPE IF EXISTS DEC_COIN;