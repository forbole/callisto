CREATE TABLE profiles_params
(
    one_row_id  BOOLEAN NOT NULL DEFAULT TRUE PRIMARY KEY,
    params      JSONB   NOT NULL,
    height      BIGINT  NOT NULL,
    CHECK (one_row_id)
);
CREATE INDEX profiles_params_height_index ON profiles_params (height);
