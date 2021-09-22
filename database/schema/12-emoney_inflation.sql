/* ---- INFLATION ---- */

CREATE TABLE emoney_inflation
(
    one_row_id              bool        PRIMARY KEY DEFAULT TRUE,
    inflation               JSONB       NOT NULL,
    last_applied_time       TIMESTAMP   NOT NULL,
    last_applied_height     INT         NOT NULL,
    CHECK (one_row_id)
);
