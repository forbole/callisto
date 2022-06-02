CREATE TABLE shield_pool
(
	pool_id				INT				NOT NULL PRIMARY KEY,
	from_address        TEXT            NOT NULL,
	shield              COIN[]          NOT NULL,
	native_deposit      COIN[]		    NOT NULL,
	foreign_deposit     COIN[]		    NOT NULL,
	sponsor             TEXT            NOT NULL,
	sponsor_address     TEXT            NOT NULL,
	description         TEXT            NOT NULL,
	shield_limit        TEXT            NOT NULL,
	pause				BOOLEAN			NOT NULL,
    height              BIGINT          NOT NULL
);
CREATE INDEX shield_pool_height_index ON shield_pool (height);
