-- +goose Down

ALTER TABLE pipelines DROP CONSTRAINT IF EXISTS pipelines_tenant_id_name_key;
ALTER TABLE caches    DROP CONSTRAINT IF EXISTS caches_tenant_id_name_key;
ALTER TABLE queues    DROP CONSTRAINT IF EXISTS queues_tenant_id_name_key;

ALTER TABLE pipelines ADD CONSTRAINT pipelines_user_id_name_key UNIQUE (user_id, name);
ALTER TABLE caches    ADD CONSTRAINT caches_user_id_name_key    UNIQUE (user_id, name);
ALTER TABLE queues    ADD CONSTRAINT queues_user_id_name_key    UNIQUE (user_id, name);
