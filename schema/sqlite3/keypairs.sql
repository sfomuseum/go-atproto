CREATE TABLE keypairs (
       did TEXT,
       label TEXT,       
       public TEXT,
       private TEXT,       
       created INTEGER,
       lastmodified INTEGER
);

CREATE UNIQUE INDEX `keypairs_by_did` ON keypairs (`did`, `label`);
CREATE INDEX `keypairs_by_created` ON keypairs (`created`);
