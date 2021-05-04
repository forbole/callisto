CREATE TABLE inflation
(
    one_row_id bool PRIMARY KEY DEFAULT TRUE,
    value     DECIMAL NOT NULL,
    height    BIGINT  NOT NULL,
    CONSTRAINT one_row_uni CHECK (one_row_id)
);
CREATE INDEX inflation_height_index ON inflation (height);