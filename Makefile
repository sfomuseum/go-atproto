SQLITE_DB = pds.db

sqlite-db:
	sqlite3 $(SQLITE_DB) < schema/sqlite3/accounts.sql
	sqlite3 $(SQLITE_DB) < schema/sqlite3/keys.sql
	sqlite3 $(SQLITE_DB) < schema/sqlite3/operations.sql
