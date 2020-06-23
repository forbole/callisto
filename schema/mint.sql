CREATE TABLE inflation
(
    value     DECIMAL                     NOT NULL,
    height    BIGINT                      NOT NULL,
    timestamp TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    PRIMARY KEY (value, height)
)
