CREATE TABLE consensus
(
    height BIGINT NOT NULL,
    round  INT    NOT NULL,
    step   TEXT   NOT NULL
);


CREATE TABLE average_block_time_per_minute 
(
  average_time DECIMAL NOT NULL, 
  timestamp TIMESTAMP WITHOUT TIMEZONE NOT NULL
);

CREATE TABLE average_block_time_per_hour
(
  average_time DECIMAL NOT NULL, 
  timestamp TIMESTAMP WITHOUT TIMEZONE NOT NULL
);

CREATE TABLE average_block_time_per_day
(
  average_time DECIMAL NOT NULL, 
  timestamp TIMESTAMP WITHOUT TIMEZONE NOT NULL
};