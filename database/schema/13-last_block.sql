-- +migrate Up

CREATE TABLE IF NOT EXISTS last_block
(
    block BIGINT NOT NULL
);


-- +migrate Down
DROP TABLE IF EXISTS last_block;
