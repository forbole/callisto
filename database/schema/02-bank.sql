-- +migrate Up
/* ---- SUPPLY ---- */

CREATE TABLE supply
(
    one_row_id BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    coins      COIN[]  NOT NULL,
    height     BIGINT  NOT NULL,
    CHECK (one_row_id)
);
CREATE INDEX supply_height_index ON supply (height);

-- +migrate Down
DROP INDEX IF EXISTS supply_height_index;
DROP TABLE IF EXISTS supply CASCADE;