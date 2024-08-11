DROP TABLE IF EXISTS posts, authors;

CREATE TABLE authors (
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE posts (
    id SERIAL PRIMARY KEY,
    author_id INTEGER REFERENCES authors(id) NOT NULL,
    title TEXT  NOT NULL,
    content TEXT NOT NULL,
    created_at BIGINT DEFAULT extract(epoch from now()),
    published_at BIGINT DEFAULT 0
);

INSERT INTO authors (name) VALUES ('Дмитрий');
INSERT INTO posts (author_id, title, content) VALUES (1 'Статья', 'Содержание статьи');