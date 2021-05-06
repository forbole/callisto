CREATE TABLE delegation_history
(
    validator_address TEXT   NOT NULL REFERENCES validator (consensus_address),
    delegator_address TEXT   NOT NULL REFERENCES account (address),
    amount            COIN   NOT NULL,
    height            BIGINT NOT NULL,
    CONSTRAINT delegation_validator_delegator_unique UNIQUE (validator_address, delegator_address, height)
);
CREATE INDEX delegation_validator_address_index ON delegation_history (validator_address);
CREATE INDEX delegation_delegator_address ON delegation_history (delegator_address);
CREATE INDEX delegation_height_index ON delegation_history (height);

CREATE TABLE redelegation_history
(
    delegator_address     TEXT                        NOT NULL REFERENCES account (address),
    src_validator_address TEXT                        NOT NULL REFERENCES validator (consensus_address),
    dst_validator_address TEXT                        NOT NULL REFERENCES validator (consensus_address),
    amount                COIN                        NOT NULL,
    completion_time       TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    height                BIGINT                      NOT NULL,
    CONSTRAINT redelegation_validator_delegator_unique UNIQUE (delegator_address, src_validator_address,
                                                               dst_validator_address, height)
);
CREATE INDEX redelegation_delegator_address_index ON redelegation_history (delegator_address);
CREATE INDEX redelegation_src_validator_address_index ON redelegation_history (src_validator_address);
CREATE INDEX redelegation_dst_validator_address_index ON redelegation_history (dst_validator_address);

CREATE TABLE unbonding_delegation_history
(
    validator_address    TEXT                        NOT NULL REFERENCES validator (consensus_address),
    delegator_address    TEXT                        NOT NULL REFERENCES account (address),
    amount               COIN                        NOT NUll,
    completion_timestamp TIMESTAMP WITHOUT TIME ZONE NOT NULL,
    height               BIGINT                      NOT NULL,
    CONSTRAINT unbonding_delegation_validator_delegator_unique UNIQUE (delegator_address, validator_address, height)
);
CREATE INDEX unbonding_delegation_validator_address_index ON unbonding_delegation_history (validator_address);
CREATE INDEX unbonding_delegation_delegator_address_index ON unbonding_delegation_history (delegator_address);