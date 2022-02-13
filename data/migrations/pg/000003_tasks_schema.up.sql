BEGIN TRANSACTION;
create table IF NOT EXISTS tasks
(
    id           uuid         not null primary key,
    dropbox_path varchar(255) not null,
    count_images int       default 0,
    started_at   timestamp default now(),
    eventid      integer
        CONSTRAINT pictures_eventid_fkey
         REFERENCES events
        ON UPDATE CASCADE ON DELETE CASCADE
);

ALTER TABLE pictures
    ADD task_id uuid
        constraint pictures_taskid_fkey
            references tasks
            on delete cascade;

ALTER TABLE tasks ADD COLUMN  eventid INT;
ALTER TABLE tasks
    ADD CONSTRAINT tasks_events_id_fk
        FOREIGN KEY (eventid) REFERENCES events
            ON UPDATE CASCADE ON DELETE CASCADE;
COMMIT;