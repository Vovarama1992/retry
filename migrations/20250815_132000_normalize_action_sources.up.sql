BEGIN;
UPDATE actions
SET source = CASE
    WHEN source IS NULL OR source = '' OR lower(source) = 'direct' THEN 'direct'
    WHEN source ~* '^utm:(tg|telegram)$' THEN 'telegram'
    WHEN source ~* '^utm:(vk|vkontakte)$' THEN 'vk'
    WHEN source ~* '^utm:(youtube|video)$' THEN 'youtube'
    WHEN source ~* '^utm:(ya|yandex)$' THEN 'yandex'
    WHEN source ~* '^utm:(google|g)$' THEN 'google'
    WHEN source ~* '^utm:social$' THEN 'social'
    WHEN source ~* '^utm:[a-z0-9_\-]+' THEN lower(regexp_replace(source, '^utm:', ''))
    WHEN source ~* '^ref:.*google\.' THEN 'google'
    WHEN source ~* '^ref:.*(ya\.ru|yandex\.)' THEN 'yandex'
    WHEN source ~* '^ref:.*(vk\.com|away\.vk\.com)' THEN 'vk'
    WHEN source ~* '^ref:.*(t\.me|telegram\.org|org\.telegram\.messenger)' THEN 'telegram'
    WHEN source ~* '^ref:.*(youtube\.|youtu\.be)' THEN 'youtube'
    ELSE 'other'
END;
CREATE INDEX IF NOT EXISTS idx_actions_source ON actions (source);
COMMIT;