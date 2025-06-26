CREATE TYPE type AS ENUM('news', 'promotion');

CREATE TABLE IF NOT EXISTS banner (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    text_uz TEXT NOT NULL,
    text_ru TEXT NOT NULL,
    text_en TEXT NOT NULL,
    title_uz VARCHAR(255) NOT NULL,
    title_ru VARCHAR(255) NOT NULL,
    title_en VARCHAR(255) NOT NULL,
    label_uz VARCHAR(255) NOT NULL,
    label_ru VARCHAR(255) NOT NULL,
    label_en VARCHAR(255) NOT NULL,
    date VARCHAR(10) NOT NULL,
    img_url TEXT,
    file_link TEXT,
    href_name VARCHAR(50),
    type type DEFAULT 'news',
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at BIGINT NOT NULL DEFAULT 0
);