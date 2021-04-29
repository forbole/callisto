CREATE TABLE genesis
(
    chain_id TEXT      NOT NULL,
    time     TIMESTAMP NOT NULL
);

CREATE TABLE consensus
(
    height BIGINT NOT NULL,
    round  INT    NOT NULL,
    step   TEXT   NOT NULL
);
CREATE INDEX consensus_height_index ON consensus (height);

CREATE TABLE average_block_time_per_minute
(
    average_time DECIMAL NOT NULL
);

CREATE TABLE average_block_time_per_hour
(
    average_time DECIMAL NOT NULL
);

CREATE TABLE average_block_time_per_day
(
    average_time DECIMAL NOT NULL
);

CREATE TABLE average_block_time_from_genesis
(
    average_time DECIMAL NOT NULL
);
