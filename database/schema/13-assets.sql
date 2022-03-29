-- +migrate Up
CREATE TABLE IF NOT EXISTS vipcoin_chain_assets_assets (
    issuer          TEXT        NOT NULL,					-- assets issuer
    name            TEXT        NOT NULL,					-- assets name
    policies        JSONB,									-- assets policies
    state           INT,									-- assets state
    issued          BIGINT,									-- assets issued
    burned          BIGINT,									-- assets burned
    withdrawn       BIGINT,									-- assets withdrawn
    in_circulation  BIGINT,									-- assets in_circulation
    properties      JSONB,									-- assets properties
    extras          JSONB									-- assets extras
);

CREATE TABLE IF NOT EXISTS vipcoin_chain_assets_set_extra (
    creator         TEXT        NOT NULL,                  -- set extra creator
    name		    TEXT        NOT NULL,                  -- set extra name
    extras          JSONB                                  -- set extra extras
);

CREATE TABLE IF NOT EXISTS vipcoin_chain_assets_manage (
   creator     	    TEXT        NOT NULL,                  -- manage assets creator
   name		    	TEXT        NOT NULL,                  -- manage assets extra name
   policies    	    JSONB,								   -- manage assets policies
   state       	    INT,								   -- manage assets state
   properties  	    JSONB,								   -- manage assets properties
   issued         	BIGINT,								   -- manage assets issued
   burned         	BIGINT,								   -- manage assets burned
   withdrawn      	BIGINT,								   -- manage assets withdrawn
   in_circulation 	BIGINT								   -- manage assets in_circulation
);

CREATE TABLE IF NOT EXISTS vipcoin_chain_assets_create (
   creator     	    TEXT        NOT NULL,                   -- create assets creator
   name			    TEXT        NOT NULL,                   -- create assets extra name
   issuer         	TEXT        NOT NULL,					-- create assets issuer
   policies    	    JSONB,									-- create assets policies
   state          	INT,									-- create assets state
   properties  	    JSONB,									-- create assets properties
   extras         	JSONB									-- create assets extras
);

-- +migrate Down
DROP TABLE IF EXISTS vipcoin_chain_assets_assets CASCADE;
DROP TABLE IF EXISTS vipcoin_chain_assets_set_extra CASCADE;
DROP TABLE IF EXISTS vipcoin_chain_assets_manage CASCADE;
DROP TABLE IF EXISTS vipcoin_chain_assets_create CASCADE;