-- +goose Up
create table feeds (
    id UUID primary key,
    created_at timestamp not null,
    updated_at timestamp not null,
    name text unique not null,
    url text not null,
    user_id UUID not null,
    CONSTRAINT fk_user_id
        FOREIGN KEY (user_id)
        REFERENCES users(id)
        ON DELETE CASCADE
);

-- +goose Down
drop table feeds;
