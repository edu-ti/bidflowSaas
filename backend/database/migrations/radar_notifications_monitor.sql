-- migrations/000007_radar_editais.up.sql
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE TABLE radar_editais (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    external_id TEXT NOT NULL,
    source TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT NOT NULL,
    value NUMERIC NOT NULL,
    deadline TIMESTAMP NOT NULL,
    raw_json JSONB NOT NULL,
    tenant_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP NULL,
    UNIQUE (external_id, source, tenant_id)
);
CREATE INDEX idx_radar_editais_tenant ON radar_editais (tenant_id);

-- migrations/000007_radar_editais.down.sql
DROP TABLE IF EXISTS radar_editais;

-- migrations/000008_radar_preferences.up.sql
CREATE TABLE radar_preferences (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    name TEXT,
    keywords TEXT[],
    min_value NUMERIC,
    categories TEXT[],
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP NULL,
    UNIQUE (tenant_id, name)
);

-- migrations/000008_radar_preferences.down.sql
DROP TABLE IF EXISTS radar_preferences;

-- migrations/000009_radar_followers.up.sql
CREATE TABLE radar_followers (
    tenant_id UUID NOT NULL,
    edital_id UUID NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP NULL,
    PRIMARY KEY (tenant_id, edital_id)
);
CREATE INDEX idx_radar_followers_tenant_edital ON radar_followers (tenant_id, edital_id);

-- migrations/000009_radar_followers.down.sql
DROP TABLE IF EXISTS radar_followers;

-- migrations/000010_notifications.up.sql
CREATE TABLE notifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    type TEXT NOT NULL,
    title TEXT NOT NULL,
    message TEXT NOT NULL,
    data JSONB,
    action_url TEXT,
    read BOOLEAN NOT NULL DEFAULT false,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP NULL
);
CREATE INDEX idx_notifications_tenant ON notifications (tenant_id);

-- migrations/000010_notifications.down.sql
DROP TABLE IF EXISTS notifications;

-- migrations/000011_monitor_snapshots.up.sql
CREATE TABLE monitor_snapshots (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    edital_id UUID NOT NULL,
    snapshot JSONB NOT NULL,
    updated_at TIMESTAMP NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP NULL
);

-- migrations/000011_monitor_snapshots.down.sql
DROP TABLE IF EXISTS monitor_snapshots;

-- migrations/000012_monitor_events.up.sql
CREATE TABLE monitor_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    tenant_id UUID NOT NULL,
    edital_id UUID NOT NULL,
    event_type TEXT NOT NULL,
    data JSONB,
    created_at TIMESTAMP NOT NULL DEFAULT now(),
    deleted_at TIMESTAMP NULL
);
CREATE INDEX idx_monitor_events_tenant ON monitor_events (tenant_id);

-- migrations/000012_monitor_events.down.sql
DROP TABLE IF EXISTS monitor_events;
