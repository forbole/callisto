/* ---- PARAMS ---- */

CREATE TABLE mint_params
(
    one_row_id            BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    mint_denom            TEXT    NOT NULL,
    inflation_rate_change DECIMAL NOT NULL,
    inflation_min         DECIMAL NOT NULL,
    inflation_max         DECIMAL NOT NULL,
    goal_bonded           DECIMAL NOT NULL,
    blocks_per_year       BIGINT  NOT NULL,
    height                BIGINT  NOT NULL,
    CHECK (one_row_id)
);

/* ---- INFLATION ---- */

CREATE TABLE inflation
(
    one_row_id bool PRIMARY KEY DEFAULT TRUE,
    value      DECIMAL NOT NULL,
    height     BIGINT  NOT NULL,
    CONSTRAINT one_row_uni CHECK (one_row_id)
);
CREATE INDEX inflation_height_index ON inflation (height);