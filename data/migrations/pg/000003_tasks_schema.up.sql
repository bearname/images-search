BEGIN TRANSACTION;
create table IF NOT EXISTS tasks
(
    id           uuid         not null primary key,
    dropbox_path varchar(255) not null,
    count_images int default 0,
    started_at timestamp default now()
);

ALTER TABLE pictures
    ADD task_id uuid
        constraint pictures_taskid_fkey
            references tasks
            on delete cascade;
COMMIT;