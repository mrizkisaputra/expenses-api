CREATE TABLE users
(
    id            varchar(100) not null,
    name          varchar(100) not null,
    email         varchar(100) not null,
    password      varchar(100) not null,
    registered_at bigint       not null
);

ALTER TABLE users
    ADD CONSTRAINT users_pk_id PRIMARY KEY (id);

ALTER TABLE users
    ADD CONSTRAINT users_email_unique UNIQUE (email);