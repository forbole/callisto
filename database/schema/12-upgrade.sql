CREATE TABLE software_upgrade_plan
(
    proposal_id     INTEGER REFERENCES proposal (id) UNIQUE,
    plan_name       TEXT        NOT NULL,
    upgrade_height  BIGINT      NOT NULL,
    info            TEXT        NOT NULL,
    height          BIGINT      NOT NULL
);
CREATE INDEX software_upgrade_plan_proposal_id_index ON software_upgrade_plan (proposal_id);
CREATE INDEX software_upgrade_plan_height_index ON software_upgrade_plan (height);
