-- +goose Up
-- Recompute users.followers_count and users.following_count from the actual
-- user_follows rows.
--
-- Background: migration 002 created DB triggers (trigger_update_follow_counts_*)
-- that maintain these counts automatically. Until the accompanying code fix, the
-- application ALSO incremented/decremented these columns inside the Follow /
-- Unfollow transactions, on top of the triggers. Every follow therefore added 2
-- and every unfollow subtracted 2, leaving the stored counts inflated (~2x).
--
-- This one-time backfill repairs the stored values so they equal the real number
-- of non-soft-deleted follow rows. It is idempotent and safe to re-run.
--
-- Note: view_count / like_count / bookmark_count are NOT touched here — their
-- app-level counterparts were either no-ops or absent, so those columns were
-- only ever maintained by their respective triggers and are already correct.

UPDATE users u
SET following_count = COALESCE(c.cnt, 0)
FROM (
    SELECT u2.id, COUNT(uf.follower_id)::bigint AS cnt
    FROM users u2
    LEFT JOIN user_follows uf
      ON uf.follower_id = u2.id AND uf.deleted_at IS NULL
    GROUP BY u2.id
) c
WHERE u.id = c.id;

UPDATE users u
SET followers_count = COALESCE(c.cnt, 0)
FROM (
    SELECT u2.id, COUNT(uf.following_id)::bigint AS cnt
    FROM users u2
    LEFT JOIN user_follows uf
      ON uf.following_id = u2.id AND uf.deleted_at IS NULL
    GROUP BY u2.id
) c
WHERE u.id = c.id;

-- +goose Down
-- No-op: this migration only repairs data; the previous (inflated) values
-- cannot be reconstructed. Rolling back leaves the corrected counts in place.
