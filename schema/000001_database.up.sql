CREATE TABLE users
(
    id         SERIAL PRIMARY KEY,
    name       VARCHAR(60) NOT NULL,
    username   VARCHAR(30) NOT NULL UNIQUE,

    email      VARCHAR(50) UNIQUE,
    phone      VARCHAR(15) UNIQUE,
    birthday   DATE,
    created_at TIMESTAMP   NOT NULL
);

CREATE TABLE posts
(
    id         SERIAL PRIMARY KEY,
    author_id  INTEGER   NOT NUll,

    content    TEXT      NOT NUll,
    created_at TIMESTAMP NOT NULL,

    FOREIGN KEY (author_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE auth
(
    user_id  INTEGER UNIQUE NOT NUll,
    password VARCHAR(255)   NOT NUll,

    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);

CREATE TABLE refresh
(
    user_id INTEGER UNIQUE NOT NULL,
    token   VARCHAR(255)   NOT NUll,

    FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE
);