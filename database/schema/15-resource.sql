CREATE TABLE resource
(
    id              TEXT    NOT NULL PRIMARY KEY,
    collection_id   TEXT,
    data            JSONB,
    name            TEXT,
    version         TEXT,
    resource_type   TEXT,
    also_known_as   JSONB,
    from_address    TEXT,
    height          BIGINT  NOT NULL
);
CREATE INDEX resource_id_index ON resource (id);
CREATE INDEX resource_height_index ON resource (height);
