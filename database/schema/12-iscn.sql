/* ---- PARAMS ---- */

CREATE TABLE iscn_params
(
    one_row_id BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    params     JSONB   NOT NULL,
    height     BIGINT  NOT NULL,
    CHECK (one_row_id)
);

/* ---- RECORD ---- */

CREATE TABLE iscn_record
(
    one_row_id  bool PRIMARY KEY DEFAULT TRUE,
    records JSONB NOT NULL,
    owner_address string NOT NULL,
    latest_version BIGINT NOT NULL,
    height      BIGINT  NOT NULL,
    CONSTRAINT one_row_uni CHECK (one_row_id)
);
CREATE INDEX iscn_record_height_index ON iscn_record (height);