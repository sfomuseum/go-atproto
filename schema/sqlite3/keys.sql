DROP TABLE IF exists keys;

CREATE TABLE keys (
       did TEXT,
       label TEXT,       
       private TEXT,       
       created INTEGER,
       lastmodified INTEGER
);

CREATE UNIQUE INDEX `keys_by_did` ON keys (`did`, `label`);
CREATE INDEX `keys_by_created` ON keys (`created`);
