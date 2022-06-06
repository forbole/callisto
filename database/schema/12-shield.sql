CREATE TABLE shield_pool
(
	pool_id				INT				NOT NULL PRIMARY KEY,
	from_address        TEXT            NOT NULL REFERENCES account (address),
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


CREATE TABLE shield_purchase
(
	pool_id				INT				NOT NULL REFERENCES shield_pool (pool_id),
	from_address        TEXT            NOT NULL PRIMARY KEY REFERENCES account (address),
	shield              COIN[]          NOT NULL,
	description         TEXT            NOT NULL,
    height              BIGINT          NOT NULL
);
CREATE INDEX shield_purchase_height_index ON shield_purchase (height);


CREATE TABLE shield_provider
(
	address        		TEXT            NOT NULL PRIMARY KEY REFERENCES account (address),
    collateral			BIGINT          NOT NULL,
    delegation_bonded	BIGINT          NOT NULL,
    native_rewards		COIN[]          NOT NULL,
    foreign_rewards		COIN[]          NOT NULL,
    total_locked		BIGINT          NOT NULL,
    withdrawing			BIGINT          NOT NULL,
    height              BIGINT          NOT NULL
);
CREATE INDEX shield_provider_height_index ON shield_purchase (height);

/* ---- PARAMS ---- */

CREATE TABLE shield_pool_params
(
    one_row_id BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    params     JSONB   NOT NULL,
    height     BIGINT  NOT NULL,
    CHECK (one_row_id)
);
CREATE INDEX shield_pool_params_height_index ON shield_pool_params (height);