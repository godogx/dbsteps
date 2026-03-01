CREATE TABLE users
(
    id         INTEGER PRIMARY KEY,
    name       TEXT NOT NULL UNIQUE,
    email      TEXT NOT NULL,
    age        INT  NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL
);

CREATE TABLE orders
(
    id         INTEGER PRIMARY KEY,
    items      INT NOT NULL,
    amount     INT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP DEFAULT NULL
);