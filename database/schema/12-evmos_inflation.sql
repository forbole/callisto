/* ---- PARAMS ---- */

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
    circulating_supply      DEC_COIN    NOT NULL,
    epochMintProvision      DEC_COIN    NOT NULL,
    inflationRate           DECIMAL     NOT NULL,
    inflationPeriod         INTEGER     NOT NULL,
    skippedEpochs           INTEGER     NOT NULL,
    height                  BIGINT      NOT NULL,
    CHECK (one_row_id)
);
