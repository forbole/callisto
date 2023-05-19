CREATE TABLE ccv_validator
(
    consumer_consensus_address TEXT NOT NULL PRIMARY KEY, /* Validator consensus address on consumer chain */
    provider_consensus_address TEXT NOT NULL UNIQUE /* Validator consensus address on provider chain */
);