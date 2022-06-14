CREATE TABLE shield_pool
(
    pool_id                 INT             NOT NULL PRIMARY KEY,
    shield                  TEXT            NOT NULL,
    native_service_fees     COIN[]          NOT NULL,
    foreign_service_fees    COIN[]          NOT NULL,
    sponsor                 TEXT            NOT NULL,
    sponsor_address         TEXT            NOT NULL,
    description             TEXT            NOT NULL,
    shield_limit            TEXT            NOT NULL,
    pause                   BOOLEAN         NOT NULL,
    height                  BIGINT          NOT NULL
);
CREATE INDEX shield_pool_height_index ON shield_pool (height);

CREATE TABLE shield_purchase
(
    pool_id				INT				NOT NULL REFERENCES shield_pool (pool_id),
    purchaser			TEXT			NOT NULL REFERENCES account (address),
    shield              TEXT			NOT NULL,
    description			TEXT			NOT NULL,
    height				BIGINT			NOT NULL
);
CREATE INDEX shield_purchase_height_index ON shield_purchase (height);

CREATE TABLE shield_provider
(
    address				TEXT			NOT NULL PRIMARY KEY REFERENCES account (address),
    collateral			BIGINT			NOT NULL,
    delegation_bonded	BIGINT			NOT NULL,
    native_rewards		COIN[]			NOT NULL,
    foreign_rewards		COIN[]			NOT NULL,
    total_locked		BIGINT			NOT NULL,
    withdrawing			BIGINT			NOT NULL,
    height				BIGINT			NOT NULL
);
CREATE INDEX shield_provider_height_index ON shield_purchase (height);

CREATE TABLE shield_withdraws
(
    address			TEXT						NOT NULL REFERENCES account (address),
    amount			BIGINT						NOT NULL,
    completion_time	TIMESTAMP WITHOUT TIME ZONE	NOT NULL,
    height			BIGINT						NOT NULL
);
CREATE INDEX shield_withdraws_height_index ON shield_withdraws (height);

CREATE TABLE shield_status
(
     one_row_id						BOOLEAN		NOT NULL DEFAULT TRUE PRIMARY KEY,
    global_staking_pool				BIGINT		NOT NULL,
    current_native_service_fees		COIN[]		NOT NULL,
    current_foreign_service_fees	COIN[]		NOT NULL,
    remaining_native_service_fees	COIN[]		NOT NULL,
    remaining_foreign_service_fees	COIN[]		NOT NULL,
    total_collateral				BIGINT		NOT NULL,
    total_shield					BIGINT		NOT NULL,
    total_withdrawing				BIGINT		NOT NULL,
    height							BIGINT		NOT NULL,
     CHECK (one_row_id)
);
CREATE INDEX shield_status_height_index ON shield_status (height);

/* ---- PARAMS ---- */

CREATE TABLE shield_pool_params
(
    one_row_id	BOOLEAN	NOT NULL DEFAULT TRUE PRIMARY KEY,
    params		JSONB	NOT NULL,
    height		BIGINT	NOT NULL,
    CHECK (one_row_id)
);
CREATE INDEX shield_pool_params_height_index ON shield_pool_params (height);

CREATE TABLE shield_claim_proposal_params
(
     one_row_id 		BOOLEAN	NOT NULL DEFAULT TRUE PRIMARY KEY,
      params			JSONB	NOT NULL,
    height			BIGINT	NOT NULL,
     CHECK (one_row_id)
);
CREATE INDEX shield_claim_proposal_params_height_index ON shield_claim_proposal_params (height);
