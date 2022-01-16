BEGIN TRANSACTION;
ALTER TABLE pictures
    DROP is_original_saved;
ALTER TABLE pictures
    DROP is_preview_saved;
ALTER TABLE pictures
    DROP is_text_recognized;
ALTER TABLE pictures
    DROP is_mobile_saved;
ALTER TABLE pictures
    DROP processing_status;
ALTER TABLE pictures
    DROP attempts;
ALTER TABLE pictures
    DROP execute_after;
COMMIT;