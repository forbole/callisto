/* ---- PARAMS ---- */

CREATE TABLE margin_params
(
    one_row_id BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    params     JSONB   NOT NULL,
    height     BIGINT  NOT NULL,
    CHECK (one_row_id)
);

CREATE TABLE margin_events
(
    transaction_hash            TEXT   NOT NULL,
    index                       BIGINT NOT NULL,
    type                        TEXT   NOT NULL,
    value                       JSONB  NOT NULL,
    height                      BIGINT NOT NULL REFERENCES block(height)
);
CREATE INDEX margin_events_type_index ON margin_events (type);
CREATE INDEX margin_events_transaction_hash_index ON margin_events (transaction_hash);
