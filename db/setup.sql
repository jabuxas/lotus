DROP TABLE IF EXISTS user;
CREATE TABLE user (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL,
    exp INTEGER
);

INSERT INTO user
    (name, exp)
VALUES
    ('jabuxas', 300),
    ('test', 10000);
