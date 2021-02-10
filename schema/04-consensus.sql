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
    average_time DECIMAL NOT NULL,
    height       BIGINT  NOT NULL UNIQUE PRIMARY KEY
);
CREATE INDEX average_block_time_per_minute_height_index ON average_block_time_per_minute (height);

CREATE TABLE average_block_time_per_hour
(
    average_time DECIMAL NOT NULL,
    height       BIGINT  NOT NULL UNIQUE PRIMARY KEY
);
CREATE INDEX average_block_time_per_hour_height_index ON average_block_time_per_hour (height);

CREATE TABLE average_block_time_per_day
(
    average_time DECIMAL NOT NULL,
    height       BIGINT  NOT NULL UNIQUE PRIMARY KEY
);
CREATE INDEX average_block_time_per_day_height_index ON average_block_time_per_day (height);

CREATE TABLE average_block_time_from_genesis
(
    average_time DECIMAL NOT NULL,
    height       BIGINT  NOT NULL UNIQUE PRIMARY KEY
);
CREATE INDEX average_block_time_from_genesis_height_index ON average_block_time_from_genesis (height);
