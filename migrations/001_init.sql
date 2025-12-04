-- 初始化数据库结构
-- 包含用户表和应用版本表

-- 用户表
CREATE TABLE IF NOT EXISTS users (
    id BIGINT PRIMARY KEY,
    openid VARCHAR(64) NOT NULL,
    nick_name VARCHAR(64),
    avatar_url VARCHAR(512),
    last_login_time BIGINT,
    created_at BIGINT DEFAULT 0,
    updated_at BIGINT DEFAULT 0,
    deleted_at BIGINT DEFAULT 0
);

CREATE UNIQUE INDEX idx_users_openid ON users(openid);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);

COMMENT ON TABLE users IS '用户信息表';
COMMENT ON COLUMN users.id IS '雪花ID主键';
COMMENT ON COLUMN users.openid IS '微信OpenID';
COMMENT ON COLUMN users.nick_name IS '昵称';
COMMENT ON COLUMN users.avatar_url IS '头像URL';
COMMENT ON COLUMN users.last_login_time IS '最后登录时间(毫秒)';
COMMENT ON COLUMN users.created_at IS '创建时间(毫秒)';
COMMENT ON COLUMN users.updated_at IS '更新时间(毫秒)';
COMMENT ON COLUMN users.deleted_at IS '软删除时间(毫秒)';

-- 应用版本表
CREATE TABLE IF NOT EXISTS app_versions (
    id BIGSERIAL PRIMARY KEY,
    version VARCHAR(20) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL DEFAULT 'Backend Template',
    description TEXT,
    min_version VARCHAR(20),
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    force_update BOOLEAN NOT NULL DEFAULT FALSE,
    release_notes TEXT,
    build_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_app_versions_active ON app_versions(is_active) WHERE is_active = TRUE;
CREATE INDEX idx_app_versions_version ON app_versions(version);

COMMENT ON TABLE app_versions IS '应用版本信息表';

-- 创建自动更新 updated_at 的触发器
CREATE OR REPLACE FUNCTION update_app_version_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

DROP TRIGGER IF EXISTS update_app_version_updated_at ON app_versions;
CREATE TRIGGER update_app_version_updated_at
    BEFORE UPDATE ON app_versions
    FOR EACH ROW
    EXECUTE FUNCTION update_app_version_updated_at();

-- 插入初始版本数据
INSERT INTO app_versions (version, name, description, is_active, build_time)
VALUES (
    '1.0.0',
    'Backend Template',
    'Initial Release',
    TRUE,
    CURRENT_TIMESTAMP
) ON CONFLICT (version) DO NOTHING;
