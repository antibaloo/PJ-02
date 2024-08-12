DROP TABLE IF EXISTS posts, authors;

CREATE TABLE authors (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    title TEXT  NOT NULL,
    content TEXT NOT NULL,
    author_id INTEGER REFERENCES authors(id) NOT NULL,
    created_at BIGINT DEFAULT extract(epoch from now()),
    published_at BIGINT DEFAULT 0
);