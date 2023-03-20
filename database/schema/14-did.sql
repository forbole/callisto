CREATE TABLE did_doc
(
    id                     TEXT   NOT NULL PRIMARY KEY,
    context                TEXT[],
    controller             TEXT[],
    verification_method    JSONB,
    authentication         TEXT[],
    assertion_method       TEXT[],
    capability_invocation  TEXT[],
    capability_delegation  TEXT[],
    key_agreement          TEXT[],
    service                JSONB,
    also_known_as          TEXT[],
    version_id             TEXT,
    height                 BIGINT  NOT NULL
);
CREATE INDEX did_doc_id_index ON did_doc (id);
CREATE INDEX did_doc_height_index ON did_doc (height);
