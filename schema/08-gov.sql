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
    yes          INTEGER   NOT NULL,
    abstain      INTEGER   NOT NULL,
    no           INTEGER   NOT NULL,
    no_with_veto INTEGER   NOT NULL,
    height       INTEGER   NOT NULL,
    timestamp    timestamp NOT NULL,
    PRIMARY KEY (proposal_id, height)
);

CREATE TABLE vote
(
    proposal_id INTEGER REFERENCES proposal (proposal_id) NOT NULL,
    voter       TEXT REFERENCES account (address),
    option      TEXT                                      NOT NULL,
    height      INTEGER                                   NOT NULL,
    timestamp   TIMESTAMP                                 NOT NULL,
    PRIMARY KEY (proposal_id, voter, height)
);

CREATE TABLE deposit
(
    proposal_id   INTEGER REFERENCES proposal (proposal_id) NOT NULL,
    depositor     TEXT REFERENCES account (address),
    amount        COIN[],
    total_deposit COIN[],
    height        INTEGER,
    timestamp     TIMESTAMP,
    PRIMARY KEY (proposal_id, depositor, height)
);
