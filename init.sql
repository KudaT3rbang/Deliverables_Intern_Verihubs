create table users
(
    id       integer generated always as identity,
    email    varchar(100) not null,
    password varchar(100) not null,
    primary key (id),
    unique (email)
);

create table books
(
    id             integer generated always as identity,
    title          varchar(100)                        not null,
    author         varchar(100)                        not null,
    published_date date                                not null,
    language       varchar(100)                        not null,
    added_at       timestamp default CURRENT_TIMESTAMP not null,
    added_by       integer                             not null,
    deleted_at     timestamp,
    deleted_by     integer,
    primary key (id),
    constraint fk_books_added_by
        foreign key (added_by) references users,
    constraint fk_books_deleted_by
        foreign key (deleted_by) references users
);

create table borrow_history
(
    id          integer generated always as identity,
    book_id     integer                             not null,
    user_id     integer                             not null,
    borrowed_at timestamp default CURRENT_TIMESTAMP not null,
    returned_at timestamp,
    primary key (id),
    constraint fk_borrow_book
        foreign key (book_id) references books,
    constraint fk_borrow_user
        foreign key (user_id) references users
);


