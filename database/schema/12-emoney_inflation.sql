/* ---- INFLATION ---- */

CREATE TABLE emoney_inflation
(
    one_row_id bool PRIMARY KEY DEFAULT TRUE,
    issuer          TEXT NOT NULL,
    denom           TEXT NOT NULL,
    rate            BIGINT NOT NULL,
    height          BIGINT  NOT NULL,
    CONSTRAINT one_row_uni CHECK (one_row_id)
);
CREATE INDEX inflation_height_index ON inflation (height);