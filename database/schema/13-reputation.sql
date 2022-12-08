CREATE TABLE reputation_feedback (
    index TEXT PRIMARY KEY NOT NULL,
    network TEXT NOT NULL,
    cptPositive INTEGER NOT NULL,
    cptNegative INTEGER NOT NULL,
    positive JSON[] NOT NULL,
    neutral JSON[] NOT NULL,
    negative JSON[] NOT NULL,
    feedbackers JSON NOT NULL,
);