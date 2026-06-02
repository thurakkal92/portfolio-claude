-- Locales reference
CREATE TABLE locales (
    code        TEXT PRIMARY KEY,
    name        TEXT NOT NULL,
    is_default  BOOLEAN NOT NULL DEFAULT FALSE
);

-- Per-locale site metadata
CREATE TABLE site_settings (
    locale            TEXT PRIMARY KEY REFERENCES locales(code) ON DELETE CASCADE,
    site_title        TEXT NOT NULL,
    site_description  TEXT NOT NULL,
    og_image_path     TEXT,
    cv_path           TEXT
);

-- Hero section
CREATE TABLE hero (
    locale               TEXT PRIMARY KEY REFERENCES locales(code) ON DELETE CASCADE,
    eyebrow              TEXT NOT NULL,
    name                 TEXT NOT NULL,
    subtitle             TEXT NOT NULL,
    cta_primary_label    TEXT NOT NULL,
    cta_secondary_label  TEXT NOT NULL
);

-- About section
CREATE TABLE about (
    locale       TEXT PRIMARY KEY REFERENCES locales(code) ON DELETE CASCADE,
    heading      TEXT NOT NULL,
    body_md      TEXT NOT NULL,
    quick_facts  JSONB NOT NULL
);

-- Skills
CREATE TABLE skill_groups (
    id             BIGSERIAL PRIMARY KEY,
    slug           TEXT NOT NULL UNIQUE,
    icon           TEXT NOT NULL,
    display_order  INT  NOT NULL
);

CREATE TABLE skill_group_translations (
    group_id  BIGINT NOT NULL REFERENCES skill_groups(id) ON DELETE CASCADE,
    locale    TEXT   NOT NULL REFERENCES locales(code)    ON DELETE CASCADE,
    title     TEXT   NOT NULL,
    PRIMARY KEY (group_id, locale)
);

CREATE TABLE skill_items (
    id             BIGSERIAL PRIMARY KEY,
    group_id       BIGINT NOT NULL REFERENCES skill_groups(id) ON DELETE CASCADE,
    label          TEXT   NOT NULL,
    display_order  INT    NOT NULL
);

CREATE INDEX skill_items_group_idx ON skill_items(group_id, display_order);

-- Projects
CREATE TABLE projects (
    id             BIGSERIAL PRIMARY KEY,
    slug           TEXT NOT NULL UNIQUE,
    company        TEXT NOT NULL,
    image_path     TEXT NOT NULL,
    live_url       TEXT,
    source_url     TEXT,
    tags           TEXT[] NOT NULL DEFAULT '{}',
    display_order  INT  NOT NULL
);

CREATE TABLE project_translations (
    project_id   BIGINT NOT NULL REFERENCES projects(id) ON DELETE CASCADE,
    locale       TEXT   NOT NULL REFERENCES locales(code) ON DELETE CASCADE,
    title        TEXT   NOT NULL,
    description  TEXT   NOT NULL,
    highlights   JSONB  NOT NULL,
    PRIMARY KEY (project_id, locale)
);

-- Experience timeline
CREATE TABLE experiences (
    id             BIGSERIAL PRIMARY KEY,
    company        TEXT NOT NULL,
    start_date     DATE NOT NULL,
    end_date       DATE,
    location       TEXT,
    display_order  INT  NOT NULL
);

CREATE TABLE experience_translations (
    experience_id  BIGINT NOT NULL REFERENCES experiences(id)  ON DELETE CASCADE,
    locale         TEXT   NOT NULL REFERENCES locales(code)    ON DELETE CASCADE,
    role           TEXT   NOT NULL,
    summary        TEXT   NOT NULL,
    PRIMARY KEY (experience_id, locale)
);

-- Contact section copy
CREATE TABLE contact (
    locale              TEXT PRIMARY KEY REFERENCES locales(code) ON DELETE CASCADE,
    heading             TEXT NOT NULL,
    subheading          TEXT NOT NULL,
    form_name_label     TEXT NOT NULL,
    form_email_label    TEXT NOT NULL,
    form_message_label  TEXT NOT NULL,
    form_submit_label   TEXT NOT NULL,
    success_message     TEXT NOT NULL,
    error_message       TEXT NOT NULL
);

-- Social links (locale-agnostic)
CREATE TABLE social_links (
    id             BIGSERIAL PRIMARY KEY,
    kind           TEXT NOT NULL,
    href           TEXT NOT NULL,
    display_label  TEXT NOT NULL,
    display_order  INT  NOT NULL
);

-- Contact form submissions audit log
CREATE TABLE contact_submissions (
    id                 BIGSERIAL PRIMARY KEY,
    name               TEXT NOT NULL,
    email              TEXT NOT NULL,
    message            TEXT NOT NULL,
    locale             TEXT,
    ip_address         INET,
    user_agent         TEXT,
    created_at         TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    email_sent         BOOLEAN NOT NULL DEFAULT FALSE,
    email_provider_id  TEXT
);

CREATE INDEX contact_submissions_created_idx ON contact_submissions(created_at DESC);

-- Per-IP rate limiting
CREATE TABLE contact_rate_limit (
    ip_address     INET PRIMARY KEY,
    attempt_count  INT  NOT NULL DEFAULT 0,
    window_start   TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
