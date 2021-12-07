CREATE TABLE fee_grant_allowance
(
    id                 SERIAL      NOT NULL PRIMARY KEY,
    grantee_address    TEXT        NOT NULL REFERENCES account (address),
    granter_address    TEXT        NOT NULL REFERENCES account (address),
    allowance          JSONB       NOT NULL,
    height             BIGINT      NOT NULL,
    CONSTRAINT unique_fee_grant_allowance UNIQUE(grantee, granter) 
);
CREATE INDEX fee_grant_allowance_height_index ON fee_grant_allowance (height);
