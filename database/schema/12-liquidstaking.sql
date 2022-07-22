CREATE TABLE liquid_staking_params
(
    one_row_id BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    params     JSONB   NOT NULL,
    height     BIGINT  NOT NULL,
    CHECK (one_row_id)
);

CREATE TABLE liquid_staking_state
(
    one_row_id BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    state      JSONB   NOT NULL,
    height     BIGINT  NOT NULL,
    CHECK (one_row_id)
);