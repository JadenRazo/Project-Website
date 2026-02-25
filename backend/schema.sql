-- =============================================
-- Portfolio Website Database Schema
-- Version: 2.0.0
-- Description: Complete database schema for portfolio website
-- including users, messaging, URL shortener, status monitoring, and more
-- =============================================

-- Enable necessary extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE EXTENSION IF NOT EXISTS "citext";

-- =============================================
-- CORE SCHEMA: Authentication & Users
-- =============================================

-- Users table (central to all features)
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username CITEXT UNIQUE NOT NULL,
    email CITEXT UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    avatar_url VARCHAR(500),
    bio TEXT,
    is_active BOOLEAN DEFAULT true,
    is_verified BOOLEAN DEFAULT false,
    role VARCHAR(50) DEFAULT 'user' CHECK (role IN ('user', 'admin', 'moderator')),
    last_login_at TIMESTAMP WITH TIME ZONE,
    failed_login_attempts INTEGER DEFAULT 0,
    locked_until TIMESTAMP WITH TIME ZONE,
    two_factor_enabled BOOLEAN DEFAULT false,
    two_factor_provider VARCHAR(50) CHECK (two_factor_provider IS NULL OR two_factor_provider IN ('google', 'github', 'microsoft', 'discord')),
    two_factor_provider_id VARCHAR(255),
    two_factor_backup_codes TEXT,
    totp_secret TEXT,
    totp_enabled BOOLEAN DEFAULT false,
    totp_verified BOOLEAN DEFAULT false,
    totp_backup_codes TEXT,
    totp_recovery_used INTEGER DEFAULT 0,
    totp_enabled_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Admin users table (for separate admin authentication if needed)
CREATE TABLE IF NOT EXISTS admin_users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    username VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    full_name VARCHAR(255),
    permissions JSONB DEFAULT '[]'::jsonb,
    is_super_admin BOOLEAN DEFAULT false,
    is_active BOOLEAN DEFAULT true,
    last_login_at TIMESTAMP WITH TIME ZONE,
    failed_login_attempts INTEGER DEFAULT 0,
    locked_until TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- User sessions for authentication
CREATE TABLE IF NOT EXISTS user_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token_hash VARCHAR(255) UNIQUE NOT NULL,
    ip_address INET,
    user_agent TEXT,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_activity TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Email verification tokens
CREATE TABLE IF NOT EXISTS email_verifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    verified_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Password reset tokens
CREATE TABLE IF NOT EXISTS password_resets (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) UNIQUE NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    used_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- API keys for developer access
CREATE TABLE IF NOT EXISTS api_keys (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    key_hash VARCHAR(255) UNIQUE NOT NULL,
    permissions JSONB DEFAULT '[]'::jsonb,
    rate_limit INTEGER DEFAULT 1000,
    expires_at TIMESTAMP WITH TIME ZONE,
    last_used_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    revoked_at TIMESTAMP WITH TIME ZONE
);

-- OAuth tokens for two-factor authentication
CREATE TABLE IF NOT EXISTS oauth_tokens (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    provider VARCHAR(50) NOT NULL CHECK (provider IN ('google', 'github', 'microsoft', 'discord')),
    encrypted_access_token TEXT NOT NULL,
    encrypted_refresh_token TEXT,
    token_expiry TIMESTAMP WITH TIME ZONE,
    token_family_id UUID,
    parent_token_id UUID,
    rotation_count INTEGER DEFAULT 0 CHECK (rotation_count <= 10),
    last_rotated_at TIMESTAMP WITH TIME ZONE,
    revoked_at TIMESTAMP WITH TIME ZONE,
    revocation_reason VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, provider)
);

-- Discord connections for linked roles (public, non-admin users)
CREATE TABLE IF NOT EXISTS discord_connections (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    discord_user_id VARCHAR(255) NOT NULL UNIQUE,
    discord_username VARCHAR(255),
    encrypted_access_token TEXT NOT NULL,
    encrypted_refresh_token TEXT,
    token_expiry TIMESTAMP WITH TIME ZONE,
    metadata_pushed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =============================================
-- AUDIT & LOGGING
-- =============================================

-- Audit log for tracking all important actions
CREATE TABLE IF NOT EXISTS audit_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    action VARCHAR(100) NOT NULL,
    entity_type VARCHAR(50),
    entity_id UUID,
    old_values JSONB,
    new_values JSONB,
    ip_address INET,
    user_agent TEXT,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- MFA events for security tracking (TOTP two-factor authentication)
CREATE TABLE IF NOT EXISTS mfa_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    event_type VARCHAR(50) NOT NULL CHECK (event_type IN ('enabled', 'disabled', 'verified', 'failed', 'recovery_used', 'backup_regenerated')),
    ip_address INET,
    user_agent TEXT,
    success BOOLEAN DEFAULT true,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- System logs
CREATE TABLE IF NOT EXISTS system_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    level VARCHAR(10) NOT NULL CHECK (level IN ('DEBUG', 'INFO', 'WARN', 'ERROR', 'FATAL')),
    service VARCHAR(50) NOT NULL,
    message TEXT NOT NULL,
    context JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =============================================
-- URL SHORTENER SCHEMA
-- =============================================

-- Shortened URLs
CREATE TABLE IF NOT EXISTS shortened_urls (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    short_code VARCHAR(20) UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    title VARCHAR(255),
    description TEXT,
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    is_active BOOLEAN DEFAULT true,
    is_private BOOLEAN DEFAULT false,
    password_hash VARCHAR(255),
    expires_at TIMESTAMP WITH TIME ZONE,
    max_clicks INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE DEFAULT NULL,
    CONSTRAINT chk_short_code_length CHECK (char_length(short_code) >= 3 AND char_length(short_code) <= 20),
    CONSTRAINT chk_original_url_length CHECK (char_length(original_url) >= 10 AND char_length(original_url) <= 2048)
);

-- URL click analytics
CREATE TABLE IF NOT EXISTS url_clicks (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    short_url_id UUID NOT NULL REFERENCES shortened_urls(id) ON DELETE CASCADE,
    clicked_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    ip_address INET,
    user_agent TEXT,
    referer TEXT,
    country_code VARCHAR(2),
    city VARCHAR(100),
    device_type VARCHAR(50),
    browser VARCHAR(50),
    os VARCHAR(50),
    is_bot BOOLEAN DEFAULT false
);

-- Alternative simplified urls table (for backwards compatibility)
CREATE TABLE IF NOT EXISTS urls (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    short_code VARCHAR(20) UNIQUE NOT NULL,
    original_url TEXT NOT NULL,
    title VARCHAR(255),
    user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    click_count INTEGER DEFAULT 0,
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- URL analytics aggregated data
CREATE TABLE IF NOT EXISTS url_analytics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    url_id UUID NOT NULL REFERENCES shortened_urls(id) ON DELETE CASCADE,
    date DATE NOT NULL,
    total_clicks INTEGER DEFAULT 0,
    unique_visitors INTEGER DEFAULT 0,
    bot_clicks INTEGER DEFAULT 0,
    mobile_clicks INTEGER DEFAULT 0,
    desktop_clicks INTEGER DEFAULT 0,
    referrer_data JSONB DEFAULT '{}',
    country_data JSONB DEFAULT '{}',
    browser_data JSONB DEFAULT '{}',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(url_id, date)
);

-- URL tags for categorization
CREATE TABLE IF NOT EXISTS url_tags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(50) UNIQUE NOT NULL,
    color VARCHAR(7),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Many-to-many relationship for URL tags
CREATE TABLE IF NOT EXISTS short_url_tags (
    short_url_id UUID NOT NULL REFERENCES shortened_urls(id) ON DELETE CASCADE,
    tag_id UUID NOT NULL REFERENCES url_tags(id) ON DELETE CASCADE,
    PRIMARY KEY (short_url_id, tag_id)
);

-- Custom domains for URL shortener
CREATE TABLE IF NOT EXISTS custom_domains (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    domain VARCHAR(255) UNIQUE NOT NULL,
    is_verified BOOLEAN DEFAULT false,
    verification_token VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =============================================
-- MESSAGING SCHEMA
-- =============================================

-- Message channels/rooms
CREATE TABLE IF NOT EXISTS messaging_channels (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE,
    description TEXT,
    type VARCHAR(50) DEFAULT 'public' CHECK (type IN ('public', 'private', 'direct')),
    icon_url VARCHAR(500),
    created_by UUID REFERENCES users(id) ON DELETE SET NULL,
    is_archived BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Channel members
CREATE TABLE IF NOT EXISTS messaging_channel_members (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    channel_id UUID NOT NULL REFERENCES messaging_channels(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role VARCHAR(50) DEFAULT 'member' CHECK (role IN ('owner', 'admin', 'moderator', 'member')),
    joined_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_read_at TIMESTAMP WITH TIME ZONE,
    is_muted BOOLEAN DEFAULT false,
    muted_until TIMESTAMP WITH TIME ZONE,
    UNIQUE(channel_id, user_id)
);

-- Messages
CREATE TABLE IF NOT EXISTS messaging_messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    channel_id UUID NOT NULL REFERENCES messaging_channels(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    parent_id UUID REFERENCES messaging_messages(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    type VARCHAR(50) DEFAULT 'text' CHECK (type IN ('text', 'image', 'file', 'system')),
    is_edited BOOLEAN DEFAULT false,
    edited_at TIMESTAMP WITH TIME ZONE,
    is_deleted BOOLEAN DEFAULT false,
    deleted_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_message_content_length CHECK (char_length(content) >= 1 AND char_length(content) <= 4000)
);

-- Message attachments
CREATE TABLE IF NOT EXISTS messaging_attachments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    message_id UUID NOT NULL REFERENCES messaging_messages(id) ON DELETE CASCADE,
    file_name VARCHAR(255) NOT NULL,
    file_url VARCHAR(500) NOT NULL,
    file_size INTEGER,
    mime_type VARCHAR(100),
    width INTEGER,
    height INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT chk_file_size_valid CHECK (file_size > 0 AND file_size <= 104857600) -- 100MB max
);

-- Message reactions
CREATE TABLE IF NOT EXISTS messaging_reactions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    message_id UUID NOT NULL REFERENCES messaging_messages(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    emoji VARCHAR(50) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(message_id, user_id, emoji)
);

-- Read receipts
CREATE TABLE IF NOT EXISTS messaging_read_receipts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    channel_id UUID NOT NULL REFERENCES messaging_channels(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    last_message_id UUID REFERENCES messaging_messages(id) ON DELETE SET NULL,
    read_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(channel_id, user_id)
);

-- Pinned messages
CREATE TABLE IF NOT EXISTS messaging_pinned_messages (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    channel_id UUID NOT NULL REFERENCES messaging_channels(id) ON DELETE CASCADE,
    message_id UUID NOT NULL REFERENCES messaging_messages(id) ON DELETE CASCADE,
    pinned_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    pinned_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(channel_id, message_id)
);

-- Message embeds
CREATE TABLE IF NOT EXISTS messaging_embeds (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    message_id UUID NOT NULL REFERENCES messaging_messages(id) ON DELETE CASCADE,
    url TEXT NOT NULL,
    title VARCHAR(255),
    description TEXT,
    image_url VARCHAR(500),
    provider_name VARCHAR(100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Word filters for content moderation
CREATE TABLE IF NOT EXISTS word_filters (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    word TEXT NOT NULL,
    severity VARCHAR(20) DEFAULT 'low' CHECK (severity IN ('low', 'medium', 'high')),
    action VARCHAR(20) DEFAULT 'flag' CHECK (action IN ('flag', 'block', 'shadowban')),
    created_by UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- =============================================
-- STATUS MONITORING SCHEMA
-- =============================================

-- Service incidents
CREATE TABLE IF NOT EXISTS incidents (
    id SERIAL PRIMARY KEY,
    service VARCHAR(255) NOT NULL,
    title VARCHAR(255) NOT NULL,
    description TEXT,
    status VARCHAR(50) NOT NULL CHECK (status IN ('investigating', 'identified', 'monitoring', 'resolved')),
    severity VARCHAR(50) NOT NULL CHECK (severity IN ('minor', 'major', 'critical')),
    started_at TIMESTAMP WITH TIME ZONE NOT NULL,
    resolved_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Status history for tracking service health over time
CREATE TABLE IF NOT EXISTS status_history (
    id SERIAL PRIMARY KEY,
    service VARCHAR(255) NOT NULL,
    status VARCHAR(50) NOT NULL,
    latency BIGINT,
    error TEXT,
    checked_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- =============================================
-- DEVELOPER PANEL & MONITORING
-- =============================================

-- System metrics for monitoring
CREATE TABLE IF NOT EXISTS metric_data (
    id BIGSERIAL PRIMARY KEY,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    latency DOUBLE PRECISION NOT NULL,
    endpoint VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Feature flags for gradual rollouts
CREATE TABLE IF NOT EXISTS feature_flags (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT,
    is_enabled BOOLEAN DEFAULT false,
    rollout_percentage INTEGER DEFAULT 0 CHECK (rollout_percentage >= 0 AND rollout_percentage <= 100),
    user_whitelist UUID[] DEFAULT ARRAY[]::UUID[],
    metadata JSONB DEFAULT '{}'::jsonb,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Rate limiting
CREATE TABLE IF NOT EXISTS rate_limits (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    identifier VARCHAR(255) NOT NULL,
    endpoint VARCHAR(255) NOT NULL,
    requests_count INTEGER NOT NULL DEFAULT 1,
    window_start TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(identifier, endpoint, window_start)
);

-- =============================================
-- CONTENT & PORTFOLIO
-- =============================================

-- Projects showcase
CREATE TABLE IF NOT EXISTS projects (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    content TEXT,
    featured_image VARCHAR(500),
    github_url VARCHAR(500),
    live_url VARCHAR(500),
    technologies TEXT[],
    status VARCHAR(50) DEFAULT 'active' CHECK (status IN ('active', 'archived', 'draft')),
    display_order INTEGER DEFAULT 0,
    view_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Project technologies (normalized many-to-many)
CREATE TABLE IF NOT EXISTS project_technologies (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_id UUID NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    technology_name VARCHAR(100) NOT NULL,
    category VARCHAR(50) CHECK (category IN ('frontend', 'backend', 'database', 'devops', 'language', 'framework', 'tool', 'other')),
    proficiency_level INTEGER CHECK (proficiency_level >= 0 AND proficiency_level <= 100),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(project_id, technology_name)
);

-- Blog posts or articles
CREATE TABLE IF NOT EXISTS posts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    content TEXT,
    excerpt TEXT,
    featured_image VARCHAR(500),
    author_id UUID REFERENCES users(id) ON DELETE SET NULL,
    status VARCHAR(50) DEFAULT 'draft' CHECK (status IN ('draft', 'published', 'archived')),
    published_at TIMESTAMP WITH TIME ZONE,
    tags TEXT[],
    view_count INTEGER DEFAULT 0,
    read_time_minutes INTEGER DEFAULT 1,
    is_featured BOOLEAN DEFAULT false,
    is_visible BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_posts_slug ON posts(slug);
CREATE INDEX IF NOT EXISTS idx_posts_status ON posts(status);
CREATE INDEX IF NOT EXISTS idx_posts_published_at ON posts(published_at DESC);
CREATE INDEX IF NOT EXISTS idx_posts_is_featured ON posts(is_featured) WHERE is_featured = true;
CREATE INDEX IF NOT EXISTS idx_posts_is_visible ON posts(is_visible) WHERE is_visible = true;

-- Skill categories for organizing skills
CREATE TABLE IF NOT EXISTS skill_categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    icon_url TEXT,
    color VARCHAR(7) CHECK (color IS NULL OR color ~ '^#[0-9A-Fa-f]{6}$'),
    sort_order INTEGER DEFAULT 1000,
    is_visible BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Skills table for developer panel
CREATE TABLE IF NOT EXISTS skills (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    category VARCHAR(50) NOT NULL CHECK (category IN ('frontend', 'backend', 'design', 'database', 'devops', 'language', 'framework', 'tool')),
    category_id UUID REFERENCES skill_categories(id) ON DELETE SET NULL,
    proficiency_level VARCHAR(50) NOT NULL CHECK (proficiency_level IN ('beginner', 'intermediate', 'advanced', 'expert')),
    proficiency_value INTEGER NOT NULL CHECK (proficiency_value >= 0 AND proficiency_value <= 100),
    is_featured BOOLEAN DEFAULT false,
    sort_order INTEGER DEFAULT 1000,
    icon_url TEXT,
    color VARCHAR(7) CHECK (color IS NULL OR color ~ '^#[0-9A-Fa-f]{6}$'),
    tags TEXT[] DEFAULT '{}',
    years_experience INTEGER,
    last_used_date DATE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- Contact form submissions
CREATE TABLE IF NOT EXISTS contact_submissions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    subject VARCHAR(255),
    message TEXT NOT NULL,
    ip_address INET,
    user_agent TEXT,
    is_read BOOLEAN DEFAULT false,
    is_spam BOOLEAN DEFAULT false,
    replied_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Code statistics
CREATE TABLE IF NOT EXISTS code_stats (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    project_name VARCHAR(255) NOT NULL,
    language VARCHAR(50) NOT NULL,
    lines_of_code INTEGER DEFAULT 0,
    files_count INTEGER DEFAULT 0,
    last_updated TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(project_name, language)
);

-- Project paths for code statistics
CREATE TABLE IF NOT EXISTS project_paths (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    path TEXT NOT NULL,
    description TEXT,
    exclude_patterns TEXT[] DEFAULT '{}',
    is_active BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE,
    UNIQUE(name, path)
);

-- Certification categories
CREATE TABLE IF NOT EXISTS certification_categories (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    sort_order INTEGER DEFAULT 1000,
    is_visible BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Certifications
CREATE TABLE IF NOT EXISTS certifications (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(255) NOT NULL,
    issuer VARCHAR(255) NOT NULL,
    credential_id VARCHAR(255),
    issue_date DATE NOT NULL,
    expiry_date DATE,
    verification_url TEXT,
    verification_text VARCHAR(255),
    badge_url TEXT,
    description TEXT,
    category_id UUID REFERENCES certification_categories(id) ON DELETE SET NULL,
    is_featured BOOLEAN DEFAULT false,
    is_visible BOOLEAN DEFAULT true,
    sort_order INTEGER DEFAULT 1000,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP WITH TIME ZONE
);

-- =============================================
-- INDEXES FOR PERFORMANCE
-- =============================================

-- User indexes
CREATE INDEX idx_users_email ON users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_username ON users(username) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_created_at ON users(created_at);
CREATE INDEX idx_users_role ON users(role) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_is_active ON users(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_users_two_factor_enabled ON users(two_factor_enabled) WHERE two_factor_enabled = true;
CREATE INDEX idx_users_two_factor_provider ON users(two_factor_provider) WHERE two_factor_provider IS NOT NULL;

-- Admin user indexes
CREATE INDEX idx_admin_users_email ON admin_users(email) WHERE deleted_at IS NULL;
CREATE INDEX idx_admin_users_username ON admin_users(username) WHERE deleted_at IS NULL;
CREATE INDEX idx_admin_users_is_super_admin ON admin_users(is_super_admin) WHERE deleted_at IS NULL;

-- Session indexes
CREATE INDEX idx_user_sessions_user_id ON user_sessions(user_id);
CREATE INDEX idx_user_sessions_expires_at ON user_sessions(expires_at);

-- OAuth tokens indexes
CREATE INDEX idx_oauth_tokens_user_id ON oauth_tokens(user_id);
CREATE INDEX idx_oauth_tokens_user_provider ON oauth_tokens(user_id, provider);
CREATE INDEX idx_oauth_tokens_expiry ON oauth_tokens(token_expiry);
CREATE INDEX idx_oauth_tokens_family_id ON oauth_tokens(token_family_id);
CREATE INDEX idx_oauth_tokens_parent_id ON oauth_tokens(parent_token_id);
CREATE INDEX idx_oauth_tokens_revoked ON oauth_tokens(revoked_at) WHERE revoked_at IS NOT NULL;

-- Discord connections indexes
CREATE INDEX idx_discord_connections_user_id ON discord_connections(discord_user_id);
CREATE INDEX idx_discord_connections_expiry ON discord_connections(token_expiry);

-- URL shortener indexes
CREATE INDEX idx_shortened_urls_short_code ON shortened_urls(short_code) WHERE is_active = true;
CREATE INDEX idx_shortened_urls_user_id ON shortened_urls(user_id);
CREATE INDEX idx_url_clicks_short_url_id ON url_clicks(short_url_id);
CREATE INDEX idx_url_clicks_clicked_at ON url_clicks(clicked_at);

-- Messaging indexes
CREATE INDEX idx_messaging_messages_channel_id ON messaging_messages(channel_id) WHERE is_deleted = false;
CREATE INDEX idx_messaging_messages_created_at ON messaging_messages(created_at);
CREATE INDEX idx_messaging_channel_members_user_id ON messaging_channel_members(user_id);
CREATE INDEX idx_messaging_read_receipts_user_channel ON messaging_read_receipts(user_id, channel_id);

-- Audit log indexes
CREATE INDEX idx_audit_logs_user_id ON audit_logs(user_id);
CREATE INDEX idx_audit_logs_entity ON audit_logs(entity_type, entity_id);
CREATE INDEX idx_audit_logs_created_at ON audit_logs(created_at);
CREATE INDEX idx_mfa_events_user_id ON mfa_events(user_id);
CREATE INDEX idx_mfa_events_created_at ON mfa_events(created_at DESC);
CREATE INDEX idx_mfa_events_event_type ON mfa_events(event_type);

-- Metrics indexes
CREATE INDEX idx_metric_data_timestamp ON metric_data(timestamp DESC);
CREATE INDEX idx_metric_data_endpoint_timestamp ON metric_data(endpoint, timestamp DESC);

-- Status monitoring indexes
CREATE INDEX idx_incidents_service ON incidents(service);
CREATE INDEX idx_incidents_status ON incidents(status);
CREATE INDEX idx_incidents_started_at ON incidents(started_at);
CREATE INDEX idx_status_history_service ON status_history(service);
CREATE INDEX idx_status_history_checked_at ON status_history(checked_at);

-- Skills indexes
CREATE INDEX idx_skills_category ON skills(category);
CREATE INDEX idx_skills_featured ON skills(is_featured);
CREATE INDEX idx_skills_sort_order ON skills(sort_order);

-- Projects indexes
CREATE INDEX idx_projects_slug ON projects(slug);
CREATE INDEX idx_projects_status ON projects(status);
CREATE INDEX idx_projects_display_order ON projects(display_order);

-- Project technologies indexes
CREATE INDEX idx_project_technologies_project_id ON project_technologies(project_id);
CREATE INDEX idx_project_technologies_category ON project_technologies(category);
CREATE INDEX idx_project_technologies_technology_name ON project_technologies(technology_name);

-- Skill categories indexes
CREATE INDEX idx_skill_categories_is_visible ON skill_categories(is_visible);
CREATE INDEX idx_skill_categories_sort_order ON skill_categories(sort_order);

-- URLs indexes (backwards compatibility table)
CREATE INDEX idx_urls_short_code ON urls(short_code) WHERE is_active = true;
CREATE INDEX idx_urls_user_id ON urls(user_id);
CREATE INDEX idx_urls_created_at ON urls(created_at);

-- URL analytics indexes
CREATE INDEX idx_url_analytics_url_id ON url_analytics(url_id);
CREATE INDEX idx_url_analytics_date ON url_analytics(date DESC);
CREATE INDEX idx_url_analytics_url_id_date ON url_analytics(url_id, date DESC);

-- Certification indexes
CREATE INDEX idx_certifications_category_id ON certifications(category_id);
CREATE INDEX idx_certifications_is_visible ON certifications(is_visible);
CREATE INDEX idx_certifications_is_featured ON certifications(is_featured);
CREATE INDEX idx_certifications_sort_order ON certifications(sort_order);
CREATE INDEX idx_certification_categories_is_visible ON certification_categories(is_visible);
CREATE INDEX idx_certification_categories_sort_order ON certification_categories(sort_order);

-- Contact submissions indexes
CREATE INDEX idx_contact_submissions_created_at ON contact_submissions(created_at);
CREATE INDEX idx_contact_submissions_email ON contact_submissions(email);

-- Project paths indexes
CREATE INDEX idx_project_paths_name ON project_paths(name) WHERE deleted_at IS NULL;
CREATE INDEX idx_project_paths_is_active ON project_paths(is_active) WHERE deleted_at IS NULL;
CREATE INDEX idx_project_paths_created_at ON project_paths(created_at);

-- =============================================
-- TRIGGERS FOR UPDATED_AT TIMESTAMPS
-- =============================================

-- Function to update updated_at timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- Apply trigger to tables with updated_at
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_admin_users_updated_at BEFORE UPDATE ON admin_users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_shortened_urls_updated_at BEFORE UPDATE ON shortened_urls
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_messaging_channels_updated_at BEFORE UPDATE ON messaging_channels
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_projects_updated_at BEFORE UPDATE ON projects
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_posts_updated_at BEFORE UPDATE ON posts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_feature_flags_updated_at BEFORE UPDATE ON feature_flags
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_skills_updated_at BEFORE UPDATE ON skills
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_word_filters_updated_at BEFORE UPDATE ON word_filters
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_incidents_updated_at BEFORE UPDATE ON incidents
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_custom_domains_updated_at BEFORE UPDATE ON custom_domains
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_certifications_updated_at BEFORE UPDATE ON certifications
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_certification_categories_updated_at BEFORE UPDATE ON certification_categories
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_project_paths_updated_at BEFORE UPDATE ON project_paths
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_urls_updated_at BEFORE UPDATE ON urls
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_url_analytics_updated_at BEFORE UPDATE ON url_analytics
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_skill_categories_updated_at BEFORE UPDATE ON skill_categories
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_oauth_tokens_updated_at BEFORE UPDATE ON oauth_tokens
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_discord_connections_updated_at BEFORE UPDATE ON discord_connections
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =============================================
-- INITIAL DATA & CONFIGURATION
-- =============================================

-- Insert default admin user (change password immediately!)
INSERT INTO users (username, email, password_hash, full_name, role, is_active, is_verified)
VALUES ('admin', 'admin@example.com', '$2a$10$DUMMY_HASH_CHANGE_THIS', 'Administrator', 'admin', true, true)
ON CONFLICT (username) DO NOTHING;

-- Insert default feature flags
INSERT INTO feature_flags (name, description, is_enabled) VALUES
    ('messaging', 'Enable messaging functionality', true),
    ('url_shortener', 'Enable URL shortener functionality', true),
    ('developer_panel', 'Enable developer panel access', true),
    ('user_registration', 'Allow new user registrations', true),
    ('status_page', 'Enable status monitoring page', true)
ON CONFLICT (name) DO NOTHING;

-- Insert default URL tags
INSERT INTO url_tags (name, color) VALUES
    ('personal', '#3B82F6'),
    ('work', '#10B981'),
    ('temporary', '#F59E0B'),
    ('social', '#8B5CF6')
ON CONFLICT (name) DO NOTHING;

-- Insert default certification categories
INSERT INTO certification_categories (name, description, sort_order) VALUES
    ('Cloud', 'Cloud computing certifications (AWS, Azure, GCP)', 100),
    ('Security', 'Cybersecurity and information security certifications', 200),
    ('Development', 'Software development and programming certifications', 300),
    ('DevOps', 'DevOps and automation certifications', 400),
    ('Data', 'Data science and analytics certifications', 500),
    ('Networking', 'Network administration and architecture certifications', 600),
    ('Project Management', 'Project and product management certifications', 700),
    ('Other', 'Other professional certifications', 1000)
ON CONFLICT (name) DO NOTHING;

-- Insert default project paths
INSERT INTO project_paths (name, path, description, exclude_patterns, is_active) VALUES
    ('Quiz Bot', '/quiz_bot', 'Interactive quiz bot application', '{"node_modules", "build", "dist", "logs", "*.log"}', true),
    ('Project Website', '/main/Project-Website', 'Portfolio website with React frontend and Go backend', '{"node_modules", "build", "dist", "logs", "*.log", "bin", "vendor"}', true)
ON CONFLICT (name, path) DO NOTHING;

-- =============================================
-- UTILITY VIEWS
-- =============================================

-- User activity summary
CREATE OR REPLACE VIEW user_activity_summary AS
SELECT 
    u.id,
    u.username,
    u.email,
    u.last_login_at,
    COUNT(DISTINCT s.id) as total_sessions,
    COUNT(DISTINCT su.id) as total_urls_created,
    COUNT(DISTINCT m.id) as total_messages_sent
FROM users u
LEFT JOIN user_sessions s ON u.id = s.user_id
LEFT JOIN shortened_urls su ON u.id = su.user_id
LEFT JOIN messaging_messages m ON u.id = m.user_id
WHERE u.deleted_at IS NULL
GROUP BY u.id;

-- URL analytics summary
CREATE OR REPLACE VIEW url_analytics_summary AS
SELECT 
    su.id,
    su.short_code,
    su.original_url,
    su.created_at,
    COUNT(uc.id) as total_clicks,
    COUNT(DISTINCT uc.ip_address) as unique_visitors,
    COUNT(DISTINCT uc.country_code) as countries_reached
FROM shortened_urls su
LEFT JOIN url_clicks uc ON su.id = uc.short_url_id
WHERE su.is_active = true
GROUP BY su.id;

-- Service status overview
CREATE OR REPLACE VIEW service_status_overview AS
SELECT 
    sh.service,
    sh.status as current_status,
    sh.latency as current_latency_ms,
    sh.checked_at as last_checked,
    COUNT(i.id) as active_incidents,
    ROUND(
        (COUNT(CASE WHEN sh2.status = 'operational' THEN 1 END)::numeric / 
         NULLIF(COUNT(sh2.id), 0) * 100), 
        2
    ) as uptime_percentage
FROM (
    SELECT DISTINCT ON (service) 
        service, status, latency, checked_at
    FROM status_history
    ORDER BY service, checked_at DESC
) sh
LEFT JOIN incidents i ON i.service = sh.service AND i.status != 'resolved'
LEFT JOIN status_history sh2 ON sh2.service = sh.service 
    AND sh2.checked_at >= CURRENT_TIMESTAMP - INTERVAL '30 days'
GROUP BY sh.service, sh.status, sh.latency, sh.checked_at;

-- =============================================
-- COMMENTS FOR DOCUMENTATION
-- =============================================

COMMENT ON TABLE users IS 'Core user accounts table for authentication and profile management';
COMMENT ON TABLE oauth_tokens IS 'Encrypted OAuth2 tokens for two-factor authentication';
COMMENT ON TABLE shortened_urls IS 'Stores shortened URLs with their metadata and settings';
COMMENT ON TABLE messaging_messages IS 'Chat messages for the messaging functionality';
COMMENT ON TABLE projects IS 'Portfolio projects showcase';
COMMENT ON TABLE audit_logs IS 'Comprehensive audit trail for all system actions';
COMMENT ON TABLE incidents IS 'Service incidents for status monitoring';
COMMENT ON TABLE status_history IS 'Historical service status data for uptime tracking';
COMMENT ON TABLE certifications IS 'Professional certifications and credentials showcase';
COMMENT ON TABLE certification_categories IS 'Categories for organizing certifications';

COMMENT ON COLUMN users.two_factor_enabled IS 'Whether 2FA is enabled for this user';
COMMENT ON COLUMN users.two_factor_provider IS 'OAuth2 provider used for 2FA (google/github/microsoft)';
COMMENT ON COLUMN users.two_factor_provider_id IS 'Provider-specific user ID for verification';
COMMENT ON COLUMN users.two_factor_backup_codes IS 'Encrypted JSON array of bcrypt hashed backup codes';

-- =============================================
-- VISITOR ANALYTICS SCHEMA (Privacy-Compliant)
-- =============================================

-- Anonymous visitor sessions (no PII)
CREATE TABLE IF NOT EXISTS visitor_sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    session_hash VARCHAR(64) NOT NULL, -- Hashed session identifier
    country_code VARCHAR(2), -- ISO country code only
    region VARCHAR(100), -- State/Province (optional based on consent)
    city VARCHAR(100), -- City (optional based on consent)
    timezone VARCHAR(50),
    language VARCHAR(10),
    device_type VARCHAR(50) CHECK (device_type IN ('desktop', 'mobile', 'tablet', 'other')),
    browser_family VARCHAR(50),
    os_family VARCHAR(50),
    is_bot BOOLEAN DEFAULT false,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_seen_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE DEFAULT (CURRENT_TIMESTAMP + INTERVAL '24 hours')
);

-- Page view analytics (no PII)
CREATE TABLE IF NOT EXISTS page_views (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    session_id UUID REFERENCES visitor_sessions(id) ON DELETE CASCADE,
    path VARCHAR(255) NOT NULL,
    referrer_domain VARCHAR(255), -- Domain only, no full URL
    duration_seconds INTEGER,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Privacy consent records
CREATE TABLE IF NOT EXISTS privacy_consents (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    session_hash VARCHAR(64) NOT NULL,
    consent_type VARCHAR(50) NOT NULL CHECK (consent_type IN ('analytics', 'functional', 'marketing')),
    granted BOOLEAN NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    expires_at TIMESTAMP WITH TIME ZONE DEFAULT (CURRENT_TIMESTAMP + INTERVAL '1 year')
);

-- Aggregated visitor metrics by time periods
CREATE TABLE IF NOT EXISTS visitor_metrics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    metric_date DATE NOT NULL,
    hour INTEGER CHECK (hour >= 0 AND hour <= 23), -- NULL for daily aggregates
    unique_visitors INTEGER DEFAULT 0,
    total_page_views INTEGER DEFAULT 0,
    avg_session_duration INTEGER DEFAULT 0,
    bounce_rate DECIMAL(5,2),
    new_visitors INTEGER DEFAULT 0,
    returning_visitors INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(metric_date, hour)
);

-- Daily visitor summaries for fast queries
CREATE TABLE IF NOT EXISTS visitor_daily_summary (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    summary_date DATE NOT NULL UNIQUE,
    unique_visitors INTEGER DEFAULT 0,
    total_sessions INTEGER DEFAULT 0,
    total_page_views INTEGER DEFAULT 0,
    avg_pages_per_session DECIMAL(5,2),
    avg_session_duration INTEGER,
    top_countries JSONB DEFAULT '{}', -- {"US": 150, "UK": 75, ...}
    top_pages JSONB DEFAULT '{}', -- {"/home": 200, "/about": 150, ...}
    device_breakdown JSONB DEFAULT '{}', -- {"mobile": 40, "desktop": 55, "tablet": 5}
    browser_breakdown JSONB DEFAULT '{}', -- {"chrome": 60, "firefox": 25, ...}
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Real-time visitor tracking (last 24 hours only)
CREATE TABLE IF NOT EXISTS visitor_realtime (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    session_hash VARCHAR(64) NOT NULL,
    last_activity TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    current_page VARCHAR(255),
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Visitor location aggregates (privacy-safe)
CREATE TABLE IF NOT EXISTS visitor_locations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    date DATE NOT NULL,
    country_code VARCHAR(2) NOT NULL,
    visitor_count INTEGER DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(date, country_code)
);

-- =============================================
-- VISITOR ANALYTICS INDEXES
-- =============================================

-- Session indexes
CREATE INDEX idx_visitor_sessions_created_at ON visitor_sessions(created_at);
CREATE INDEX idx_visitor_sessions_expires_at ON visitor_sessions(expires_at);
CREATE INDEX idx_visitor_sessions_hash ON visitor_sessions(session_hash);
CREATE INDEX idx_visitor_sessions_last_seen ON visitor_sessions(last_seen_at DESC);

-- Page view indexes
CREATE INDEX idx_page_views_session_id ON page_views(session_id);
CREATE INDEX idx_page_views_created_at ON page_views(created_at);
CREATE INDEX idx_page_views_path ON page_views(path);

-- Metrics indexes
CREATE INDEX idx_visitor_metrics_date_hour ON visitor_metrics(metric_date DESC, hour);
CREATE INDEX idx_visitor_daily_summary_date ON visitor_daily_summary(summary_date DESC);
CREATE INDEX idx_visitor_realtime_activity ON visitor_realtime(last_activity DESC);
CREATE INDEX idx_visitor_locations_date_country ON visitor_locations(date DESC, country_code);

-- Privacy consent indexes
CREATE INDEX idx_privacy_consents_session ON privacy_consents(session_hash);
CREATE INDEX idx_privacy_consents_expires ON privacy_consents(expires_at);

-- =============================================
-- VISITOR ANALYTICS VIEWS
-- =============================================

-- Current visitor summary view
CREATE OR REPLACE VIEW visitor_current_summary AS
SELECT 
    COUNT(DISTINCT session_hash) as active_visitors,
    COUNT(DISTINCT current_page) as active_pages,
    AVG(EXTRACT(EPOCH FROM (CURRENT_TIMESTAMP - created_at))) as avg_session_seconds
FROM visitor_realtime
WHERE last_activity >= CURRENT_TIMESTAMP - INTERVAL '5 minutes';

-- Privacy-compliant visitor overview
CREATE OR REPLACE VIEW visitor_analytics_overview AS
SELECT 
    vds.summary_date,
    vds.unique_visitors,
    vds.total_page_views,
    vds.avg_session_duration,
    vds.avg_pages_per_session,
    COALESCE(vds.top_countries, '{}'::jsonb) as top_countries,
    COALESCE(vds.top_pages, '{}'::jsonb) as top_pages
FROM visitor_daily_summary vds
ORDER BY vds.summary_date DESC;

-- =============================================
-- VISITOR ANALYTICS FUNCTIONS
-- =============================================

-- Function to clean up expired sessions
CREATE OR REPLACE FUNCTION cleanup_expired_visitor_sessions()
RETURNS void AS $$
BEGIN
    -- Delete expired sessions
    DELETE FROM visitor_sessions WHERE expires_at < CURRENT_TIMESTAMP;
    
    -- Delete old real-time data
    DELETE FROM visitor_realtime WHERE created_at < CURRENT_TIMESTAMP - INTERVAL '24 hours';
    
    -- Delete expired consents
    DELETE FROM privacy_consents WHERE expires_at < CURRENT_TIMESTAMP;
END;
$$ LANGUAGE plpgsql;

-- Function to aggregate visitor metrics
CREATE OR REPLACE FUNCTION aggregate_visitor_metrics(target_date DATE)
RETURNS void AS $$
DECLARE
    hour_num INTEGER;
BEGIN
    -- Aggregate hourly metrics
    FOR hour_num IN 0..23 LOOP
        INSERT INTO visitor_metrics (
            metric_date, hour, unique_visitors, total_page_views,
            avg_session_duration, new_visitors, returning_visitors
        )
        SELECT 
            target_date,
            hour_num,
            COUNT(DISTINCT vs.session_hash),
            COUNT(pv.id),
            AVG(EXTRACT(EPOCH FROM (vs.last_seen_at - vs.created_at))),
            COUNT(DISTINCT CASE WHEN vs.created_at::date = target_date THEN vs.session_hash END),
            COUNT(DISTINCT CASE WHEN vs.created_at::date < target_date THEN vs.session_hash END)
        FROM visitor_sessions vs
        LEFT JOIN page_views pv ON vs.id = pv.session_id
        WHERE DATE_TRUNC('hour', vs.created_at) = target_date::timestamp + (hour_num || ' hours')::interval
        ON CONFLICT (metric_date, hour) DO UPDATE SET
            unique_visitors = EXCLUDED.unique_visitors,
            total_page_views = EXCLUDED.total_page_views,
            avg_session_duration = EXCLUDED.avg_session_duration,
            new_visitors = EXCLUDED.new_visitors,
            returning_visitors = EXCLUDED.returning_visitors;
    END LOOP;
END;
$$ LANGUAGE plpgsql;

-- =============================================
-- VISITOR ANALYTICS TRIGGERS
-- =============================================

-- Trigger to update last_seen_at on page view
CREATE OR REPLACE FUNCTION update_session_last_seen()
RETURNS TRIGGER AS $$
BEGIN
    UPDATE visitor_sessions 
    SET last_seen_at = CURRENT_TIMESTAMP
    WHERE id = NEW.session_id;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_update_session_last_seen
AFTER INSERT ON page_views
FOR EACH ROW EXECUTE FUNCTION update_session_last_seen();

-- =============================================
-- SCHEDULED JOBS (pg_cron extension required)
-- =============================================

-- Schedule cleanup job (requires pg_cron)
-- SELECT cron.schedule('cleanup-visitor-sessions', '0 * * * *', 'SELECT cleanup_expired_visitor_sessions();');
-- SELECT cron.schedule('aggregate-visitor-metrics', '5 0 * * *', 'SELECT aggregate_visitor_metrics(CURRENT_DATE - INTERVAL ''1 day'');');

-- =============================================
-- COMMENTS FOR VISITOR ANALYTICS
-- =============================================

COMMENT ON TABLE visitor_sessions IS 'Anonymous visitor sessions with privacy-compliant tracking';
COMMENT ON TABLE page_views IS 'Page view events linked to visitor sessions';
COMMENT ON TABLE privacy_consents IS 'Records of user privacy consent preferences';
COMMENT ON TABLE visitor_metrics IS 'Aggregated visitor metrics by hour and day';
COMMENT ON TABLE visitor_daily_summary IS 'Pre-computed daily summaries for fast queries';
COMMENT ON TABLE visitor_realtime IS 'Real-time visitor tracking for current activity';
COMMENT ON TABLE visitor_locations IS 'Aggregated visitor counts by location';

-- =============================================
-- END OF SCHEMA
-- =============================================
-- DevPanel Prompt Tables
-- Create UUID extension if not exists
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Create prompt_categories table
CREATE TABLE IF NOT EXISTS prompt_categories (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    name VARCHAR(255) UNIQUE NOT NULL,
    description TEXT,
    is_visible BOOLEAN DEFAULT true,
    sort_order INTEGER DEFAULT 100,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create index on deleted_at for soft deletes
CREATE INDEX IF NOT EXISTS idx_prompt_categories_deleted_at ON prompt_categories(deleted_at);

-- Create prompts table
CREATE TABLE IF NOT EXISTS prompts (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    prompt TEXT NOT NULL,
    category_id UUID REFERENCES prompt_categories(id),
    is_featured BOOLEAN DEFAULT false,
    is_visible BOOLEAN DEFAULT true,
    sort_order INTEGER DEFAULT 100,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP
);

-- Create indexes
CREATE INDEX IF NOT EXISTS idx_prompts_deleted_at ON prompts(deleted_at);
CREATE INDEX IF NOT EXISTS idx_prompts_category_id ON prompts(category_id);

-- Add updated_at trigger function if not exists
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

-- Add triggers for updated_at
DROP TRIGGER IF EXISTS update_prompt_categories_updated_at ON prompt_categories;
CREATE TRIGGER update_prompt_categories_updated_at
    BEFORE UPDATE ON prompt_categories
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

DROP TRIGGER IF EXISTS update_prompts_updated_at ON prompts;
CREATE TRIGGER update_prompts_updated_at
    BEFORE UPDATE ON prompts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- =====================================================================================
-- Privacy configuration table for visitor analytics
-- =====================================================================================

CREATE TABLE IF NOT EXISTS privacy_configs (
    id VARCHAR(50) PRIMARY KEY,
    enable_tracking BOOLEAN DEFAULT true,
    privacy_mode VARCHAR(20) DEFAULT 'balanced',
    data_collection JSONB,
    retention JSONB,
    compliance JSONB,
    anonymization_options JSONB,
    consent_requirements JSONB,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_by VARCHAR(255)
);

-- Insert default config
INSERT INTO privacy_configs (id, enable_tracking, privacy_mode, data_collection, retention, compliance, anonymization_options, consent_requirements)
VALUES ('default', true, 'balanced', 
    '{"collectCookies": false, "collectIPAddresses": false, "collectUserAgents": true, "collectReferrers": true, "collectGeographicData": true, "collectSessionData": true, "collectEventData": true, "collectDeviceInfo": true, "collectBrowserInfo": true, "respectDNT": true, "anonymousMode": false}',
    '{"enableAutoDelete": true, "sessionDataDays": 30, "pageViewDataDays": 90, "aggregatedDataDays": 365, "consentRecordDays": 365, "deleteInactiveAfter": 180, "retentionPolicy": "standard"}',
    '{"gdpr": {"enabled": true, "requireConsent": false, "allowPortability": true, "allowErasure": true, "processingBasis": "legitimate_interest"}, "ccpa": {"enabled": true, "allowOptOut": true, "provideDataDisclosure": true, "doNotSellData": true}, "lgpd": {"enabled": true, "requireConsent": false, "allowPortability": true, "allowErasure": true, "dataProcessingBasis": "legitimate_interest"}, "pipeda": {"enabled": true, "requireConsent": false, "limitDataCollection": true, "provideAccess": true}}',
    '{"anonymizeIP": true, "ipAnonymizationMode": "hash", "hashSessionIDs": true, "removePII": true, "useFingerprinting": false, "maskUserAgents": false}',
    '{"requireExplicitConsent": false, "consentCategories": ["necessary", "analytics", "functional"], "defaultConsent": "opt-out", "consentDuration": 365, "showBanner": false, "bannerPosition": "bottom", "allowGranularControl": true, "minimumAge": 13}'
)
ON CONFLICT (id) DO NOTHING;

-- =============================================
-- BLOG SCHEMA
-- =============================================

CREATE TABLE IF NOT EXISTS posts (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title VARCHAR(255) NOT NULL,
    slug VARCHAR(255) UNIQUE NOT NULL,
    content TEXT,
    excerpt TEXT,
    featured_image VARCHAR(500),
    author_id UUID REFERENCES users(id) ON DELETE SET NULL,
    status VARCHAR(50) DEFAULT 'draft',
    published_at TIMESTAMP,
    tags TEXT[],
    view_count INTEGER DEFAULT 0,
    read_time_minutes INTEGER DEFAULT 1,
    is_featured BOOLEAN DEFAULT false,
    is_visible BOOLEAN DEFAULT true,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_posts_slug ON posts(slug);
CREATE INDEX IF NOT EXISTS idx_posts_status ON posts(status);
CREATE INDEX IF NOT EXISTS idx_posts_published_at ON posts(published_at DESC);
CREATE INDEX IF NOT EXISTS idx_posts_is_featured ON posts(is_featured) WHERE is_featured = true;
