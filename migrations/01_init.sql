-- Создаём таблицы 

CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash TEXT NOT NULL,
    role VARCHAR(50) NOT NULL CHECK (role IN ('owner', 'admin', 'controller', 'viewer')),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE homes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    owner_id UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE rooms (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    home_id UUID NOT NULL REFERENCES homes(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL
);

CREATE TABLE devices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    room_id UUID REFERENCES rooms(id),
    type VARCHAR(50) NOT NULL CHECK (type IN ('light', 'temperature_sensor', 'humidity_sensor', 'electricity_meter', 'water_meter', 'power_switch')),
    name VARCHAR(255) NOT NULL,
    status VARCHAR(20) NOT NULL DEFAULT 'offline' CHECK (status IN ('online', 'offline', 'disabled')),
    last_seen TIMESTAMPTZ,
    metadata JSONB,
    version INTEGER NOT NULL DEFAULT 1
);

CREATE TABLE device_states (
    device_id UUID PRIMARY KEY REFERENCES devices(id),
    state JSONB NOT NULL DEFAULT '{}',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE audit_log (
    id BIGSERIAL PRIMARY KEY,
    user_id UUID NOT NULL REFERENCES users(id),
    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(50) NOT NULL,
    entity_id VARCHAR(255) NOT NULL,
    params JSONB,
    timestamp TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    correlation_id UUID
);

CREATE INDEX idx_audit_log_timestamp ON audit_log(timestamp DESC);

CREATE TABLE telemetry (
    time TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    device_id UUID NOT NULL REFERENCES devices(id),
    value DOUBLE PRECISION NOT NULL,
    type VARCHAR(50) NOT NULL CHECK (type IN ('temperature', 'humidity', 'power', 'energy', 'water_flow'))
);

-- SELECT create_hypertable('telemetry', 'time', if_not_exists => TRUE);
CREATE INDEX idx_telemetry_device_time ON telemetry(device_id, time DESC);

CREATE TABLE scenes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    home_id UUID NOT NULL REFERENCES homes(id),
    name VARCHAR(255) NOT NULL,
    condition_json JSONB NOT NULL DEFAULT '{}',
    action_json JSONB NOT NULL DEFAULT '{}',
    is_active BOOLEAN NOT NULL DEFAULT false,
    created_by UUID NOT NULL REFERENCES users(id),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    version INTEGER NOT NULL DEFAULT 1
);

-- Добавим пару тестовых устройств для проверки

INSERT INTO users (id, username, email, password_hash, role) VALUES
('a0000000-0000-0000-0000-000000000001', 'testuser', 'test@example.com', 'hashed_password', 'owner');

INSERT INTO homes (id, name, owner_id) VALUES
('b0000000-0000-0000-0000-000000000001', 'Мой дом', 'a0000000-0000-0000-0000-000000000001');

INSERT INTO rooms (id, home_id, name) VALUES
('c0000000-0000-0000-0000-000000000001', 'b0000000-0000-0000-0000-000000000001', 'Гостиная');

INSERT INTO devices (id, room_id, type, name, status) VALUES
('d0000000-0000-0000-0000-000000000001', 'c0000000-0000-0000-0000-000000000001', 'light', 'Лампа потолочная', 'online'),
('d0000000-0000-0000-0000-000000000002', 'c0000000-0000-0000-0000-000000000001', 'temperature_sensor', 'Датчик температуры', 'online');

INSERT INTO device_states (device_id, state) VALUES
('d0000000-0000-0000-0000-000000000001', '{"power": "on", "brightness": 80}'),
('d0000000-0000-0000-0000-000000000002', '{"temperature": 22.5}');