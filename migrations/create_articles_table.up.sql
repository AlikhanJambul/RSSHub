CREATE TABLE articles (
                          id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
                          created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                          updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
                          title TEXT NOT NULL,
                          link TEXT NOT NULL,
                          published_at TIMESTAMP,
                          description TEXT,
                          feed_id UUID NOT NULL REFERENCES feeds(id) ON DELETE CASCADE
);

CREATE UNIQUE INDEX articles_feed_link_idx ON articles(feed_id, link);
