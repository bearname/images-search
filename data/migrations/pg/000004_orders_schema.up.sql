BEGIN TRANSACTION;

CREATE TABLE IF NOT EXISTS orders
(
    id             uuid   not null primary key,
    userId         int REFERENCES users ON DELETE CASCADE,
    status         smallint     default 0,
    totalPrice     bigint not null,
    updatedAt      timestamp    default now(),
    pay_id         varchar(255) default '',
    receipt_url    varchar(255) default '',
    receipt_number varchar(255) default ''
);

CREATE TABLE IF NOT EXISTS picture_in_order
(
    orderId   uuid REFERENCES orders ON DELETE CASCADE,
    pictureId uuid REFERENCES pictures ON DELETE CASCADE
);

COMMIT;