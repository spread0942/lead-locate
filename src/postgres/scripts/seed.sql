-- Create tables
CREATE TABLE users (
    username varchar(150) NOT NULL,
    password varchar(255) NOT NULL,
    email varchar(255) NOT NULL,
    date_stored timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL,
    date_update timestamp DEFAULT CURRENT_TIMESTAMP NOT NULL
);