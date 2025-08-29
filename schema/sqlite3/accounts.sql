CREATE TABLE accounts (
       did TEXT PRIMARY KEY,
       handle TEXT,
       created INTEGER,
       lastmodified INTEGER
);

CREATE UNIQUE INDEX `accounts_by_handle` ON accounts (`handle`);
CREATE INDEX `accounts_by_created` ON accounts (`created`);