CREATE TABLE evmos_inflation_params
(
    one_row_id BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    params     JSONB   NOT NULL,
    height     BIGINT  NOT NULL,
    CHECK (one_row_id)
);

CREATE TABLE evmos_inflation_data
(
    one_row_id              BOOLEAN     NOT NULL DEFAULT TRUE PRIMARY KEY,
    circulating_supply      DEC_COIN[]  NOT NULL,
    epoch_mint_provision    DEC_COIN[]  NOT NULL,
    inflation_rate          DECIMAL     NOT NULL,
    inflation_period        INTEGER     NOT NULL,
    skipped_epochs          INTEGER     NOT NULL,
    height                  BIGINT      NOT NULL,
    CHECK (one_row_id)
);
