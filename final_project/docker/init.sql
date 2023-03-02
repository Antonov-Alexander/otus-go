-- Checkers
CREATE TABLE checkers (
    id   SERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL
);

-- Limit values
CREATE TABLE limits (
    id              SERIAL PRIMARY KEY,
    checker_id      INT NOT NULL,
    item            JSON,
    interval_value  INT NOT NULL,
    limit_value     INT NOT NULL
);

-- Black/White lists
CREATE TYPE LIST_TYPE AS ENUM('white', 'black');
CREATE TABLE lists (
    id          SERIAL PRIMARY KEY,
    checker_id  INT NOT NULL,
    item        JSON NOT NULL,
    list_type   LIST_TYPE
);

-- DATA
INSERT INTO checkers VALUES
(1, 'ip'),
(2, 'login'),
(3, 'password');

INSERT INTO limits (checker_id, interval_value, limit_value) VALUES
(1, 10, 5),
(2, 10, 5),
(3, 10, 5);

INSERT INTO lists (checker_id, item, list_type) VALUES
(1, '{"ip":123456}'::JSON, 'white'),
(1, '{"ip":654321}'::JSON, 'black'),
(2, '{"login":"brute"}'::JSON, 'white'),
(2, '{"login":"force"}'::JSON, 'black'),
(3, '{"password":"pass123456"}'::JSON, 'white'),
(3, '{"password":"pass654321"}'::JSON, 'black');
