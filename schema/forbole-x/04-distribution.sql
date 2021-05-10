CREATE TYPE DEC_COIN AS
(
    denom  TEXT,
    amount TEXT
);

CREATE TABLE delegation_reward_history
(
    validator_address TEXT       NOT NULL,
    delegator_address TEXT       NOT NULL REFERENCES account (address),
    withdraw_address  TEXT       NOT NULL,
    amount            DEC_COIN[] NOT NULL,
    height            BIGINT     NOT NULL REFERENCES block(height),
    CONSTRAINT validator_delegator_unique UNIQUE (delegator_address, validator_address, height)
);
CREATE INDEX delegation_reward_delegator_address_index ON delegation_reward_history (delegator_address);
CREATE INDEX delegation_reward_height_index ON delegation_reward_history (height);

CREATE TABLE validator_commission_amount_history
(
    self_delegate_address TEXT       NOT NULL,
    amount                DEC_COIN[] NOT NULL,
    height                BIGINT     NOT NULL REFERENCES block(height),
    CONSTRAINT commission_height_unique UNIQUE (self_delegate_address, height)
);
CREATE INDEX validator_commission_amount_height_index ON validator_commission_amount_history (height);