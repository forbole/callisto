CREATE TABLE inflation
(
    value  DECIMAL NOT NULL,
    height BIGINT  NOT NULL,
    PRIMARY KEY (value, height)
);
