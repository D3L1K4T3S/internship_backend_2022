CREATE DATABASE avito;

CREATE TABLE users(
    id SERIAL PRIMARY KEY,
    balance REAL CHECK(balance >= 0),
    reserved REAL DEFAULT 0 CHECK(reserved >= 0)
);

CREATE TABLE services(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE orders(
    id SERIAL PRIMARY KEY,
    user_id INTEGER,
    service_id INTEGER,
    status TEXT DEFAULT 'Not completed' CHECK(status IN('Done','Not completed','Not enough funds')),
    cost REAL CHECK (cost >= 0),
    date TIMESTAMP DEFAULT NOW()
);

CREATE TABLE transactions(
    id SERIAL PRIMARY KEY,
    order_id INTEGER NOT NULL,
    type TEXT CHECK(type IN('Debiting','Refund')),
    time_trans TIMESTAMP DEFAULT NOW(),
    amount REAL NOT NULL,
    pass BOOLEAN DEFAULT TRUE,
    comment TEXT
);

CREATE TABLE deposits(
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL,
    time_trans TIMESTAMP DEFAULT NOW(),
    amount REAL NOT NULL,
    comment TEXT
);

ALTER TABLE orders ADD FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE orders ADD FOREIGN KEY (service_id) REFERENCES services (id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE transactions ADD FOREIGN KEY (order_id) REFERENCES orders (id) ON DELETE CASCADE ON UPDATE CASCADE;

ALTER TABLE deposits ADD FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE ON UPDATE CASCADE;

INSERT INTO services(name) VALUES
    ('Delivery'),
    ('Subscription'),
    ('Guarantor'),
    ('Fitting'),
    ('AdditionalInformation');

