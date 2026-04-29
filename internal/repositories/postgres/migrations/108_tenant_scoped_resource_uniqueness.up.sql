-- +goose Up

-- Pipelines, caches, and queues were created with UNIQUE(user_id, name), which
-- allowed two users in the same tenant to register identically named
-- resources. The application layer scopes lookups by tenant_id, so the
-- service-layer name check sees only one row while the duplicate slips past
-- the user-scoped DB constraint. Replace with tenant-scoped uniqueness.

-- Disambiguate any preexisting duplicates within the same tenant by appending
-- a short user-id suffix to all but the oldest collider. Required so the new
-- UNIQUE(tenant_id, name) constraint can be added without conflicts.
WITH ranked AS (
    SELECT id,
           tenant_id,
           name,
           user_id,
           ROW_NUMBER() OVER (
               PARTITION BY tenant_id, name
               ORDER BY created_at, id
           ) AS rn
    FROM pipelines
    WHERE tenant_id IS NOT NULL
)
UPDATE pipelines p
SET name = p.name || '-' || substr(replace(r.user_id::text, '-', ''), 1, 8)
FROM ranked r
WHERE p.id = r.id AND r.rn > 1;

WITH ranked AS (
    SELECT id,
           tenant_id,
           name,
           user_id,
           ROW_NUMBER() OVER (
               PARTITION BY tenant_id, name
               ORDER BY created_at, id
           ) AS rn
    FROM caches
    WHERE tenant_id IS NOT NULL
)
UPDATE caches c
SET name = c.name || '-' || substr(replace(r.user_id::text, '-', ''), 1, 8)
FROM ranked r
WHERE c.id = r.id AND r.rn > 1;

WITH ranked AS (
    SELECT id,
           tenant_id,
           name,
           user_id,
           ROW_NUMBER() OVER (
               PARTITION BY tenant_id, name
               ORDER BY created_at, id
           ) AS rn
    FROM queues
    WHERE tenant_id IS NOT NULL
)
UPDATE queues q
SET name = q.name || '-' || substr(replace(r.user_id::text, '-', ''), 1, 8)
FROM ranked r
WHERE q.id = r.id AND r.rn > 1;

ALTER TABLE pipelines DROP CONSTRAINT IF EXISTS pipelines_user_id_name_key;
ALTER TABLE caches    DROP CONSTRAINT IF EXISTS caches_user_id_name_key;
ALTER TABLE queues    DROP CONSTRAINT IF EXISTS queues_user_id_name_key;

ALTER TABLE pipelines ADD CONSTRAINT pipelines_tenant_id_name_key UNIQUE (tenant_id, name);
ALTER TABLE caches    ADD CONSTRAINT caches_tenant_id_name_key    UNIQUE (tenant_id, name);
ALTER TABLE queues    ADD CONSTRAINT queues_tenant_id_name_key    UNIQUE (tenant_id, name);
