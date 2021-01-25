CREATE TABLE proposal
(
    title             TEXT      NOT NULL,
    description       TEXT      NOT NULL,
    proposal_route    TEXT      NOT NULL,
    proposal_type     TEXT      NOT NULL,
    proposal_id       DECIMAL   NOT NULL PRIMARY KEY,
    submit_time       TIMESTAMP NOT NULL,
    deposit_end_time  TIMESTAMP,
    voting_start_time TIMESTAMP,
    voting_end_time   TIMESTAMP,
    proposer          TEXT      NOT NULL REFERENCES account (address),
    status            TEXT
);

CREATE TABLE tally_result
(
    proposal_id  INTEGER REFERENCES proposal (proposal_id),
    yes          BIGINT NOT NULL,
    abstain      BIGINT NOT NULL,
    no           BIGINT NOT NULL,
    no_with_veto BIGINT NOT NULL,
    height       BIGINT NOT NULL,
    PRIMARY KEY (proposal_id, height)
);

CREATE TABLE vote
(
    proposal_id INTEGER REFERENCES proposal (proposal_id) NOT NULL,
    voter       TEXT REFERENCES account (address),
    option      TEXT                                      NOT NULL,
    height      BIGINT                                    NOT NULL,
    PRIMARY KEY (proposal_id, voter, height)
);

CREATE TABLE deposit
(
    proposal_id INTEGER REFERENCES proposal (proposal_id) NOT NULL,
    depositor   TEXT REFERENCES account (address),
    amount      COIN[],
    height      BIGINT,
    PRIMARY KEY (proposal_id, depositor, height)
);
