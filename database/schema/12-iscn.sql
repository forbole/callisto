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
    one_row_id           BOOLEAN PRIMARY KEY DEFAULT TRUE,
 	ipld                 TEXT       NOT NULL,
	context              JSONB      NOT NULL,
	record_id            TEXT       NOT NULL,
	record_type          TEXT       NOT NULL,
	content_fingerprints JSONB      NOT NULL,
	content_metadata     JSONB      NOT NULL,
	record_notes         TEXT,       
	record_timestamp     TIMESTAMP,
	record_version       BIGINT     NOT NULL,
	stakeholders         JSONB      NOT NULL,
	height               BIGINT     NOT NULL,
    CONSTRAINT one_row_uni CHECK (one_row_id)
);
CREATE INDEX iscn_record_height_index ON iscn_record (height);