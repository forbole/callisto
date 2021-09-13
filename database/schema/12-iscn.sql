/* ---- PARAMS ---- */

CREATE TABLE iscn_params
(
    one_row_id BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    params     JSONB   NOT NULL,
    height     BIGINT  NOT NULL,
    CHECK (one_row_id)
);

-- /* ---- RECORD ---- */

CREATE TABLE iscn_record
(
    iscn_id         TEXT      NOT NULL PRIMARY KEY,
    owner_address   TEXT      NOT NULL,
    latest_version  BIGINT    NOT NULL,
    ipld            TEXT      NOT NULL,
    iscn_data       JSONB     NOT NULL,
    height          BIGINT    NOT NULL
);
CREATE INDEX iscn_record_height_index ON iscn_record (height);
