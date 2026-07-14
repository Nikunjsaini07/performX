CREATE TABLE performance_review_comment_likes (
    comment_id UUID NOT NULL REFERENCES performance_review_comments(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT now(),

    PRIMARY KEY (comment_id, user_id)
);

CREATE INDEX performance_review_comment_likes_user_idx ON performance_review_comment_likes(user_id);
