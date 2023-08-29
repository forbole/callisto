CREATE TABLE ccv_validator
(
    consumer_consensus_address      TEXT   NOT NULL PRIMARY KEY, /* Validator consensus address on consumer chain */
    consumer_self_delegate_address  TEXT   NOT NULL, /* Validator self delegate address on consumer chain */
    consumer_operator_address       TEXT   NOT NULL, /* Validator operator address on consumer chain */
    provider_consensus_address      TEXT   NOT NULL, /* Validator consensus address on provider chain */
    provider_self_delegate_address  TEXT   NOT NULL, /* Validator self delegate address on provider chain */
    provider_operator_address       TEXT   NOT NULL, /* Validator operator address on provider chain */
    height                          BIGINT NOT NULL
);