CREATE TABLE marker_params
(
    one_row_id BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    params     JSONB   NOT NULL,
    height     BIGINT  NOT NULL,
    CHECK (one_row_id)
);

CREATE TABLE marker_account
(
    address                  TEXT       NOT NULL REFERENCES account (address),
    access_control           TEXT       NOT NULL,
    allow_governance_control BOOLEAN    NOT NULL,
    denom                    TEXT       NOT NULL PRIMARY KEY,
    marker_type              TEXT       NOT NULL,
    status                   TEXT       NOT NULL,
    supply                   TEXT       NOT NULL,
    height                   BIGINT     NOT NULL
);
CREATE INDEX marker_account_height_index ON marker_account (height);
