CREATE TABLE IF NOT EXISTS PaymentTypes (
    id SERIAL PRIMARY KEY,
    title VARCHAR(30) UNIQUE
);

CREATE TABLE IF NOT EXISTS ServiceTypes (
    id SERIAL PRIMARY KEY,
    title VARCHAR(30) UNIQUE
);

CREATE TABLE IF NOT EXISTS StatusTypes (
    id serial primary key,
    title VARCHAR(10) UNIQUE
);

CREATE TABLE IF NOT EXISTS PayeeNames (
    id serial primary key,
    title VARCHAR(10) UNIQUE
);

CREATE TABLE IF NOT EXISTS PaymentNarratives (
    id serial primary key,
    title VARCHAR(10) UNIQUE
);

CREATE TABLE IF NOT EXISTS Transcations (
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
    statusid INT,
    paymenttypeid INT,
    paymentnumber VARCHAR(10),
    serviceid INT,
    servicetypeid INT,
    payeeid INT,
    payeenameid INT,
    payeebankmfo INT,
    payeebankaccount VARCHAR(17),
    paymentnarrativeid int,
    foreign key (statusid) references StatusTypes(id),
    foreign key (paymenttypeid) references PaymentTypes(id),
    foreign key (servicetypeid) references ServiceTypes(id),
    foreign key (payeenameid) references PayeeNames(id),
    foreign key (paymentnarrativeid) references PaymentNarratives(id)
);