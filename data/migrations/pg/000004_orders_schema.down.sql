BEGIN TRANSACTION;

ALTER TABLE picture_in_order DROP CONSTRAINT picture_in_order.picture_in_order_orderid_fkey;
ALTER TABLE picture_in_order DROP CONSTRAINT picture_in_order.picture_in_order_pictureid_fkey;
TRUNCATE TABLE orders;
TRUNCATE TABLE picture_in_order;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS picture_in_order;

COMMIT;