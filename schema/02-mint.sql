CREATE TABLE inflation_history
(
    value  DECIMAL NOT NULL,
    height BIGINT  NOT NULL,
    PRIMARY KEY (value, height)
);
