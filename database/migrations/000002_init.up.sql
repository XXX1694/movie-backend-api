create table if not exists users (
                                     id serial primary key,
                                     name varchar(255) not null,
    email varchar(255),
    age int,
    created_at timestamp default now(),
    deleted_at timestamp
    );

create table if not exists audit_log (
                                         id serial primary key,
                                         user_id int,
                                         action varchar(255),
    created_at timestamp default now()
    );

insert into users (name, email, age) values ('John Doe', 'john@example.com', 25);