CREATE TABLE fee_grant_allowance
(
    id                 SERIAL      NOT NULL PRIMARY KEY,
    grantee            TEXT        NOT NULL REFERENCES account (address),
    granter            TEXT        NOT NULL REFERENCES account (address),
    allowance          JSONB       NOT NULL,
    height             BIGINT      NOT NULL REFERENCES block (height)
);
CREATE INDEX fee_grant_allowance_height_index ON fee_grant_allowance (height);
