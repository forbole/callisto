-- +migrate Up
CREATE TABLE modules
(
    module_name TEXT NOT NULL UNIQUE PRIMARY KEY
);

-- +migrate Down
DROP TABLE IF EXISTS modules CASCADE;