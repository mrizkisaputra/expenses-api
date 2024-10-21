-- mengatifkan fitur fungsi uuid_generate_v4()
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE EXTENSION IF NOT EXISTS CITEXT;

CREATE TABLE users
(
    id           UUID DEFAULT uuid_generate_v4(),
    first_name   VARCHAR(100) NOT NULL,
    last_name    VARCHAR(100) NOT NULL,
    email        VARCHAR(100) NOT NULL UNIQUE,
    password     VARCHAR(100) NOT NULL,
    avatar       VARCHAR(512),
    city         VARCHAR(100),
    phone_number VARCHAR(13),
    created_at   BIGINT       NOT NULL,
    updated_at   BIGINT       NOT NULL
);

CREATE TABLE expenses
(
    id          UUID                     DEFAULT uuid_generate_v4(),
    id_user     UUID           NOT NULL,
    description VARCHAR(200)   NOT NULL,
    amount      DECIMAL(10, 2) NOT NULL,
    category    VARCHAR(100)   NOT NULL,
    created_at  BIGINT         NOT NULL,
    updated_at  BIGINT         NOT NULL,
    deleted_at  TIMESTAMP WITH TIME ZONE DEFAULT NULL
);

-------------------------------------------- CONSTRAINT TYPE CHECK --------------------------------------------
ALTER TABLE users
    ADD CONSTRAINT first_name_check CHECK ( users.first_name <> '' ),
    ADD CONSTRAINT last_name_check CHECK ( users.last_name <> '' ),
    ADD CONSTRAINT email_check CHECK ( users.email <> '' ),
    ADD CONSTRAINT password_check CHECK ( users.password <> '' );

ALTER TABLE expenses
    ADD CONSTRAINT description_check CHECK ( expenses.description <> '' ),
    ADD CONSTRAINT category_check CHECK ( expenses.category <> '' ),
    ADD CONSTRAINT amount_check CHECK ( expenses.amount > 0 );

-------------------------------------------- CONSTRAINT PRIMARY KEY --------------------------------------------
ALTER TABLE users
    ADD CONSTRAINT users_id_pk PRIMARY KEY (id);

ALTER TABLE expenses
    ADD CONSTRAINT expenses_id_pk PRIMARY KEY (id);

-------------------------------------------- CONSTRAINT FOREIGN KEY --------------------------------------------
ALTER TABLE expenses
    ADD CONSTRAINT expenses_id_user_fk FOREIGN KEY (id_user) REFERENCES users (id)
        ON DELETE CASCADE ON UPDATE RESTRICT;

---------------------------------------------------- INDEX ------------------------------------------------------
CREATE INDEX IF NOT EXISTS expenses_id_user_index ON expenses (id_user);
CREATE INDEX IF NOT EXISTS expenses_created_at_index ON expenses (created_at);
