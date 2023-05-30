CREATE TABLE pool
(
    id                           INTEGER   NOT NULL PRIMARY KEY,
    name                         TEXT      NOT NULL UNIQUE,
    runtime                      TEXT      NOT NULL,
    logo                         TEXT      NOT NULL,
    config                       TEXT      NOT NULL,
    start_key                    TEXT      NOT NULL,
    current_key                  TEXT      NOT NULL,
    current_summary              TEXT      NOT NULL,
    current_index                INTEGER   NOT NULL,
    total_bundles                INTEGER   NOT NULL,
    upload_interval              INTEGER   NOT NULL,
    operating_cost               INTEGER   NOT NULL,
    min_delegation               INTEGER   NOT NULL,
    max_bundle_size              INTEGER   NOT NULL,
    disabled                     BOOLEAN   NOT NULL,
    funders                      JSONB     NOT NULL,
    total_funds                  INTEGER   NOT NULL,
    protocol                     JSONB     NOT NULL,
    upgrade_plan                 JSONB     NOT NULL,
    current_storage_provider_id  INTEGER   NOT NULL,
    current_compression_id       INTEGER NOT NULL,
    height                       BIGINT  NOT NULL
);