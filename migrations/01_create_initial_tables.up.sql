CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS CITEXT;

CREATE TYPE role AS ENUM ('admin', 'user');
CREATE TYPE status AS ENUM ('pending', 'accepted');

DROP TABLE IF EXISTS users CASCADE;
CREATE TABLE users
(
    user_id    UUID PRIMARY KEY                  DEFAULT uuid_generate_v4(),
    first_name VARCHAR(32)              NOT NULL CHECK ( first_name <> '' ),
    last_name  VARCHAR(32)              NOT NULL CHECK ( last_name <> '' ),
    email      VARCHAR(64) UNIQUE       NOT NULL CHECK ( email <> '' ),
    avatar     VARCHAR(250),
    password   VARCHAR(250)             NOT NULL CHECK ( octet_length(password) <> 0 ),
    role       role                     NOT NULL DEFAULT 'user',

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP WITH TIME ZONE          DEFAULT CURRENT_TIMESTAMP
);

DROP TABLE IF EXISTS librarians CASCADE;
CREATE TABLE librarians
(
    librarian_id  UUID PRIMARY KEY                  DEFAULT uuid_generate_v4(),
    first_name VARCHAR(32)              NOT NULL CHECK ( first_name <> '' ),
    last_name  VARCHAR(32)              NOT NULL CHECK ( last_name <> '' ),
    email      VARCHAR(64) UNIQUE       NOT NULL CHECK ( email <> '' ),
    avatar     VARCHAR(250),
    password   VARCHAR(250)             NOT NULL CHECK ( octet_length(password) <> 0 ),

    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

DROP TABLE IF EXISTS orders CASCADE;
CREATE TABLE orders
(
    order_id    UUID PRIMARY KEY                  DEFAULT uuid_generate_v4(),
    user_id     UUID REFERENCES users (user_id),
    librarian_id   UUID REFERENCES librarians (librarian_id),
    item        JSONB,
    status      status        NOT NULL DEFAULT 'pending',

    pickup_schedule  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at  TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);