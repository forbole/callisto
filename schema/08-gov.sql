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
    proposer_address  TEXT      NOT NULL REFERENCES account (address),
    status            TEXT
);
CREATE INDEX proposal_proposer_address_index ON proposal (proposer_address);

CREATE TABLE proposal_deposit
(
    proposal_id       INTEGER REFERENCES proposal (proposal_id) NOT NULL,
    depositor_address TEXT REFERENCES account (address),
    amount            COIN[],
    height            BIGINT,
    PRIMARY KEY (proposal_id, depositor_address, height)
);
CREATE INDEX proposal_deposit_proposal_id_index ON proposal_deposit (proposal_id);
CREATE INDEX proposal_deposit_depositor_address_index ON proposal_deposit (depositor_address);

CREATE TABLE proposal_vote
(
    proposal_id   INTEGER NOT NULL REFERENCES proposal (proposal_id),
    voter_address TEXT    NOT NULL REFERENCES account (address),
    option        TEXT    NOT NULL,
    height        BIGINT  NOT NULL,
    PRIMARY KEY (proposal_id, voter_address, height)
);
CREATE INDEX proposal_vote_proposal_id_index ON proposal_vote (proposal_id);
CREATE INDEX proposal_vote_voter_address_index ON proposal_vote (voter_address);
CREATE INDEX proposal_vote_height_index ON proposal_vote (height);

CREATE TABLE proposal_tally_result
(
    proposal_id  INTEGER REFERENCES proposal (proposal_id),
    yes          BIGINT NOT NULL,
    abstain      BIGINT NOT NULL,
    no           BIGINT NOT NULL,
    no_with_veto BIGINT NOT NULL,
    height       BIGINT NOT NULL,
    PRIMARY KEY (proposal_id, height)
);
CREATE INDEX proposal_tally_result_proposal_id_index ON proposal_tally_result (proposal_id);
CREATE INDEX proposal_tally_result_height_index ON proposal_tally_result (height);