CREATE TABLE IF NOT EXISTS ServiceTypes (
    id SERIAL PRIMARY KEY,
    title VARCHAR(30) UNIQUE
);

CREATE TYPE statuss AS ENUM ('accepted', 'declined');

CREATE TYPE paymenttype AS ENUM ('cash', 'card');

CREATE TABLE IF NOT EXISTS PayeeNames (
    id serial primary key,
    title VARCHAR(10) UNIQUE
);

CREATE TABLE IF NOT EXISTS PaymentNarratives (
    id serial primary key,
    title VARCHAR(255) UNIQUE
);

CREATE TABLE IF NOT EXISTS Transactions (
    id serial primary key,
    requestid INT,
    terminalid INT,
    partnerobjectid INT,
    amounttotal FLOAT,
    amountoriginal FLOAT,
    commisionps FLOAT,
    commisionclient FLOAT,
    commisionprovider FLOAT,
    dateinput TIMESTAMP,
    datepost TIMESTAMP,
    statusid statuss,
    paymenttype paymenttype,
    paymentnumber VARCHAR(10),
    serviceid INT,
    servicetypeid INT,
    payeeid INT,
    payeenameid INT,
    payeebankmfo INT,
    payeebankaccount VARCHAR(17),
    paymentnarrativeid int,
    foreign key (servicetypeid) references ServiceTypes(id) ON DELETE CASCADE ON UPDATE CASCADE,
    foreign key (payeenameid) references PayeeNames(id) ON DELETE CASCADE ON UPDATE CASCADE,
    foreign key (paymentnarrativeid) references PaymentNarratives(id) ON DELETE CASCADE ON UPDATE CASCADE
);