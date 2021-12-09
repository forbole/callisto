CREATE TABLE fee_grant_allowance
(
    grantee            TEXT        NOT NULL PRIMARY KEY,
    granter            TEXT        NOT NULL,
    allowance          JSONB       NOT NULL,
    height             BIGINT      NOT NULL
);
CREATE INDEX fee_grant_allowance_height_index ON fee_grant_allowance (height);
