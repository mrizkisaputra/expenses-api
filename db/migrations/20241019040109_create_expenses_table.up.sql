CREATE TABLE expenses(
    id varchar(100) not null,
    id_user varchar(100) not null,
    description varchar(500),
    amount decimal(10, 2) not null,
    category varchar(100) not null,
    created_at bigint not null,
    updated_at bigint not null,
    deleted_at timestamp default null
);

ALTER TABLE expenses
ADD CONSTRAINT expenses_pk_id PRIMARY KEY (id);

ALTER TABLE expenses
ADD CONSTRAINT expenses_fk_id_user FOREIGN KEY (id_user) REFERENCES users(id);