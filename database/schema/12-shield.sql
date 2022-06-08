CREATE TABLE shield_pool
(
	pool_id				INT				NOT NULL PRIMARY KEY,
	from_address        TEXT            NOT NULL REFERENCES account (address),
	shield              TEXT         	NOT NULL,
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

CREATE TABLE shield_purchase_list
(
	purchase_id 			INT 						NOT NULL PRIMARY KEY,
	pool_id					INT							NOT NULL REFERENCES shield_pool (pool_id),
	purchaser       		TEXT            			NOT NULL REFERENCES account (address),
	deletion_time			TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	protection_end_time 	TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	foreign_service_fees	COIN[] 						NOT NULL,
	native_service_fees 	COIN[] 						NOT NULL,
	shield          	   	TEXT         				NOT NULL,
	description         	TEXT            			NOT NULL,
    height              	BIGINT          			NOT NULL
);
CREATE INDEX shield_purchase_list_height_index ON shield_purchase_list (height);

CREATE TABLE shield_withdraws
(
	address			TEXT            			NOT NULL REFERENCES account (address),
	amount 			BIGINT 						NOT NULL,
	completion_time TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    height          BIGINT          			NOT NULL
);
CREATE INDEX shield_withdraws_height_index ON shield_withdraws (height);

CREATE TABLE shield_info
(
 	one_row_id 						BOOLEAN 					NOT NULL DEFAULT TRUE PRIMARY KEY,
	global_staking_pool 			TEXT 						NOT NULL,
	last_update_time				TIMESTAMP WITHOUT TIME ZONE NOT NULL,
	next_pool_id					BIGINT 						NOT NULL,
	next_purchase_id				BIGINT						NOT NULL,
	original_staking 				JSONB 						NOT NULL DEFAULT '{}'::JSONB,
	proposal_id_reimbursement_pair	JSONB 						NOT NULL DEFAULT '{}'::JSONB,
	shield_admin					TEXT 						NOT NULL,
	shield_staking_rate				TEXT 						NOT NULL,
	stake_for_shields				JSONB 						NOT NULL DEFAULT '{}'::JSONB,
	total_claimed					BIGINT 						NOT NULL,
	total_collateral				BIGINT 						NOT NULL,
	total_shield					BIGINT 						NOT NULL,
	total_withdrawing				BIGINT 						NOT NULL,
    height              			BIGINT          			NOT NULL,
 	CHECK (one_row_id)
);
CREATE INDEX shield_info_height_index ON shield_info (height);

CREATE TABLE shield_service_fees
(
    one_row_id 						BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    foreign_service_fees     		COIN [] NOT NULL,
    native_service_fees     		COIN [] NOT NULL,
    remaining_foreign_service_fees  COIN [] NOT NULL,
	remaining_native_service_fees 	COIN [] NOT NULL,
    height     						BIGINT  NOT NULL,
    CHECK (one_row_id)
);
CREATE INDEX shield_service_fees_height_index ON shield_service_fees (height);

/* ---- PARAMS ---- */

CREATE TABLE shield_pool_params
(
    one_row_id BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    params     JSONB   NOT NULL,
    height     BIGINT  NOT NULL,
    CHECK (one_row_id)
);
CREATE INDEX shield_pool_params_height_index ON shield_pool_params (height);

CREATE TABLE shield_claim_proposal_params
(
 	one_row_id 		BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
  	params 			JSONB 	NOT NULL,
    height          BIGINT	NOT NULL,
 	CHECK (one_row_id)
);
CREATE INDEX shield_claim_proposal_params_height_index ON shield_claim_proposal_params (height);