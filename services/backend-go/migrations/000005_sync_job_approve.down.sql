-- Remove 'approve' dos tipos válidos de sync_job
ALTER TABLE public.sync_job DROP CONSTRAINT IF EXISTS sync_job_type_check;
ALTER TABLE public.sync_job ADD CONSTRAINT sync_job_type_check
    CHECK (type IN ('sync_uc','calculate','generate_pdf','recalculate_cycle'));
