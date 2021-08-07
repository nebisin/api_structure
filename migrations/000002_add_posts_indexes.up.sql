CREATE INDEX IF NOT EXISTS posts_title_idx ON posts USING gin (to_tsvector('simple', title));
CREATE INDEX IF NOT EXISTS posts_tags_idx ON posts USING gin (tags);
