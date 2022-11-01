CREATE INDEX IF NOT EXISTS tids ON transactions (terminalid);

CREATE INDEX IF NOT EXISTS dpost ON transactions (datepost);

CREATE INDEX IF NOT EXISTS stat ON transactions (statusid);

CREATE INDEX IF NOT EXISTS ptype ON transactions (paymenttype);


