CREATE TYPE PROVIDER_INFO AS
(
    email       TEXT,
    website     TEXT
);

CREATE TABLE akash_provider
(
    owner_address   TEXT            NOT NULL REFERENCES account (address),
    host_uri        TEXT            NOT NULL,
    attributes      JSONB           NOT NULL DEFAULT '[]'::JSONB,
    info            PROVIDER_INFO   NOT NULL,
    height          BIGINT          NOT NULL,
    CONSTRAINT unique_provider UNIQUE (owner_address)
);
CREATE INDEX provider_address_index ON akash_provider (owner_address);

CREATE TYPE AKASH_RESOURCE AS
(
    cpu                 BIGINT,
    memory              BIGINT,
    storage_ephemeral   BIGINT
);

CREATE TABLE akash_provider_inventory
(
    provider_address            TEXT            NOT NULL REFERENCES akash_provider (owner_address),
    active                      BOOLEAN         NOT NULL,
    lease_count                 BIGINT          NOT NULL,
    bidengine_order_count       BIGINT          NOT NULL,
    manifest_deployment_count   BIGINT          NOT NULL,
    cluster_public_hostname     TEXT            NOT NULL,
    inventory_status_raw        JSONB           NOT NULL DEFAULT '{}'::JSONB,
    active_inventory_sum        AKASH_RESOURCE  NOT NULL,
    pending_inventory_sum       AKASH_RESOURCE  NOT NULL,
    available_inventory_sum     AKASH_RESOURCE  NOT NULL,
    height                      BIGINT          NOT NULL,
    CONSTRAINT unique_provider_address UNIQUE (provider_address)
);
CREATE INDEX provider_inventory_address_index ON akash_provider_inventory (provider_address);
