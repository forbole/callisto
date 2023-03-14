/* ---- CONFIG ---- */

CREATE TABLE wormhole_config
(
    one_row_id BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    config     JSONB   NOT NULL,
    height     BIGINT  NOT NULL,
    CHECK (one_row_id)
);
CREATE INDEX wormhole_config_height_index ON wormhole_config (height);