ALTER TABLE domain_events_v2
    ADD COLUMN IF NOT EXISTS source_copy_id FixedString(32),
    MODIFY ORDER BY (scope_type, scope_key, drain_cursor, idempotency_key, source_copy_id),
    ADD CONSTRAINT IF NOT EXISTS source_copy_id_present
        CHECK source_copy_id != toFixedString('', 32)
;

ALTER TABLE domain_events_v2 ADD CONSTRAINT source_copy_id_present CHECK source_copy_id != toFixedString('', 32);

ALTER TABLE domain_events_v2 ADD CONSTRAINT IF NOT EXISTS source_copy_id_present CHECK source_copy_id != toFixedString('', 32);
