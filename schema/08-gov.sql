CREATE TABLE proposal(
    title TEXT NOT NULL,
	description TEXT NOT NULL,
	proposal_route TEXT NOT NULL,
	proposal_type TEXT NOT NULL,
	proposal_ID DECIMAL NOT NULL,
	status TEXT NOT NULL, 
	submit_time TIMESTAMP,
	deposit_end_time TIMESTAMP,
	total_deposit COIN,
	voting_start_time TIMESTAMP,
	voting_end_time TIMESTAMP,
    PRIMARY KEY (proposalID)
);

CREATE TABLE tally_result(
    ProposalID INTEGER REFERENCES proposal (proposalID),
    Yes        INTEGER,
    Abstain    INTEGER,
    No         INTEGER,
    NoWithVeto INTEGER,
    Height     INTEGER,
    Timestamp  timestamp,
    PRIMARY KEY (ProposalID,timestamp)
);
