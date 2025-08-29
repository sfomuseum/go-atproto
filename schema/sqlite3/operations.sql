CREATE TABLE operations (
       cid TEXT PRIMARY KEY,
       did TEXT,
       operation TEXT,
       created INTEGER,
       lastmodified INTEGER
);

CREATE INDEX `operations_by_did` ON operations (`did`);
CREATE INDEX `operations_by_created` ON operations (`created`);