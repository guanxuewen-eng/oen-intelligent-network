-- OEN Intelligent Network Database Migration
-- Created: 2026-04-13
-- Description: Initial schema for all 10 core tables

-- 1. Agent table
CREATE TABLE IF NOT EXISTS agent (
    id              SERIAL PRIMARY KEY,
    agent_key       VARCHAR(128) NOT NULL UNIQUE,
    name            VARCHAR(256) NOT NULL,
    role            VARCHAR(128),
    state           VARCHAR(64) DEFAULT 'inactive',
    route_mode      VARCHAR(64),
    last_heartbeat  TIMESTAMP WITH TIME ZONE,
    metadata        JSONB,
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- 2. Agent heartbeat table
CREATE TABLE IF NOT EXISTS agent_heartbeat (
    id              SERIAL PRIMARY KEY,
    agent_id        INTEGER NOT NULL REFERENCES agent(id),
    heartbeat_at    TIMESTAMP WITH TIME ZONE NOT NULL,
    status          VARCHAR(64) NOT NULL,
    route_mode      VARCHAR(64),
    cpu_usage       DOUBLE PRECISION,
    memory_usage    DOUBLE PRECISION,
    error_message   TEXT
);

CREATE INDEX IF NOT EXISTS idx_agent_heartbeat_agent_id ON agent_heartbeat(agent_id);
CREATE INDEX IF NOT EXISTS idx_agent_heartbeat_heartbeat_at ON agent_heartbeat(heartbeat_at);

-- 3. Artifact table
CREATE TABLE IF NOT EXISTS artifact (
    id                  SERIAL PRIMARY KEY,
    artifact_key        VARCHAR(128) NOT NULL UNIQUE,
    artifact_type       VARCHAR(128) NOT NULL,
    title               VARCHAR(512) NOT NULL,
    description         TEXT,
    target_system       VARCHAR(256),
    applicable_version  VARCHAR(128),
    risk_level          VARCHAR(64),
    verification_status VARCHAR(64),
    creator_agent_id    INTEGER REFERENCES agent(id),
    source_draft_id     INTEGER,
    current_version_id  INTEGER,
    metadata            JSONB,
    created_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    updated_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

-- 4. Artifact version table
CREATE TABLE IF NOT EXISTS artifact_version (
    id              SERIAL PRIMARY KEY,
    artifact_id     INTEGER NOT NULL REFERENCES artifact(id),
    version_number  VARCHAR(64) NOT NULL,
    change_summary  TEXT,
    content_json    JSONB,
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    created_by      VARCHAR(128)
);

CREATE INDEX IF NOT EXISTS idx_artifact_version_artifact_id ON artifact_version(artifact_id);

-- 5. Artifact view table
CREATE TABLE IF NOT EXISTS artifact_view (
    id                   SERIAL PRIMARY KEY,
    artifact_version_id  INTEGER NOT NULL REFERENCES artifact_version(id),
    view_type            VARCHAR(64) NOT NULL,
    content              JSONB,
    created_at           TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_artifact_view_version_id ON artifact_view(artifact_version_id);

-- 6. Candidate resource table
CREATE TABLE IF NOT EXISTS candidate_resource (
    id              SERIAL PRIMARY KEY,
    candidate_key   VARCHAR(128) NOT NULL UNIQUE,
    agent_id        INTEGER REFERENCES agent(id),
    source_type     VARCHAR(64),
    title           VARCHAR(512) NOT NULL,
    summary         TEXT,
    raw_content     TEXT,
    state           VARCHAR(64) DEFAULT 'pending',
    artifact_id     INTEGER REFERENCES artifact(id),
    review_notes    TEXT,
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    reviewed_at     TIMESTAMP WITH TIME ZONE,
    reviewed_by     VARCHAR(128)
);

-- 7. Recommendation table
CREATE TABLE IF NOT EXISTS recommendation (
    id                   SERIAL PRIMARY KEY,
    recommendation_key   VARCHAR(128) NOT NULL UNIQUE,
    agent_id             INTEGER NOT NULL REFERENCES agent(id),
    artifact_id          INTEGER NOT NULL REFERENCES artifact(id),
    title                VARCHAR(512) NOT NULL,
    match_reason         TEXT,
    environment_hint     TEXT,
    risk_summary         TEXT,
    confidence_score     DOUBLE PRECISION,
    suggested_action     TEXT,
    state                VARCHAR(64) DEFAULT 'pending',
    created_at           TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    decided_at           TIMESTAMP WITH TIME ZONE
);

CREATE INDEX IF NOT EXISTS idx_recommendation_agent_id ON recommendation(agent_id);
CREATE INDEX IF NOT EXISTS idx_recommendation_artifact_id ON recommendation(artifact_id);

-- 8. Recommendation decision table
CREATE TABLE IF NOT EXISTS recommendation_decision (
    id                  SERIAL PRIMARY KEY,
    recommendation_id   INTEGER NOT NULL UNIQUE REFERENCES recommendation(id),
    decision            VARCHAR(64) NOT NULL,
    decided_by          VARCHAR(128) NOT NULL,
    decided_at          TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    note                TEXT
);

-- 9. Consent record table
CREATE TABLE IF NOT EXISTS consent_record (
    id              SERIAL PRIMARY KEY,
    agent_id        INTEGER NOT NULL REFERENCES agent(id),
    consent_type    VARCHAR(128) NOT NULL,
    status          VARCHAR(64) DEFAULT 'active',
    granted_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW(),
    revoked_at      TIMESTAMP WITH TIME ZONE,
    granted_by      VARCHAR(128)
);

CREATE INDEX IF NOT EXISTS idx_consent_record_agent_id ON consent_record(agent_id);

-- 10. Audit log table
CREATE TABLE IF NOT EXISTS audit_log (
    id              SERIAL PRIMARY KEY,
    event_type      VARCHAR(128) NOT NULL,
    event_subtype   VARCHAR(128),
    actor_type      VARCHAR(64) NOT NULL,
    actor_id        VARCHAR(128),
    target_type     VARCHAR(64),
    target_id       VARCHAR(128),
    description     TEXT,
    metadata        JSONB,
    created_at      TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_audit_log_event_type ON audit_log(event_type);
CREATE INDEX IF NOT EXISTS idx_audit_log_created_at ON audit_log(created_at);
