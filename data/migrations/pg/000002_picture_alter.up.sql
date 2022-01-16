BEGIN TRANSACTION;
ALTER TABLE pictures
    ADD dropbox_path varchar default null;
ALTER TABLE pictures
    ADD is_original_saved bool default false;
ALTER TABLE pictures
    ADD is_preview_saved bool default false;
ALTER TABLE pictures
    ADD is_text_recognized bool default false;
ALTER TABLE pictures
    ADD is_mobile_saved bool default false;
ALTER TABLE pictures
    ADD processing_status smallint default null;
ALTER TABLE pictures
    ADD attempts smallint default 0;
ALTER TABLE pictures
    ADD execute_after timestamp DEFAULT now();
ALTER TABLE pictures
    ADD update_at timestamp DEFAULT now();
COMMIT;