CREATE TABLE marker_params
(
    one_row_id BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    params     JSONB   NOT NULL,
    height     BIGINT  NOT NULL,
    CHECK (one_row_id)
);

CREATE TABLE marker_account
(
    id                       SERIAL     NOT NULL PRIMARY KEY,
    access_control           JSONB      NOT NULL DEFAULT '{}'::JSONB,
    allow_governance_control BOOLEAN    NOT NULL,
    base_account             JSONB      NOT NULL DEFAULT '{}'::JSONB,
    denom                    TEXT       NOT NULL,
    marker_type              TEXT       NOT NULL,
    status                   TEXT       NOT NULL,
    supply                   TEXT       NOT NULL,
    height                   BIGINT     NOT NULL
);
CREATE INDEX marker_account_height_index ON marker_account (height);


CREATE TABLE marker_acc
(
    id         SERIAL  NOT NULL PRIMARY KEY,
    marker     JSONB   NOT NULL,
    height     BIGINT  NOT NULL
);
CREATE INDEX marker_acc_height_index ON marker_acc (height);
