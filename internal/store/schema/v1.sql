PRAGMA user_version = 1;

CREATE TABLE user(
  id INTEGER PRIMARY KEY,
  username TEXT UNIQUE NOT NULL,
  hashed_password BLOB NOT NULL,
  session_token TEXT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE user_post(
  id INTEGER PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES user(id) ON DELETE CASCADE,
  title TEXT NOT NULL,
  url TEXT NOT NULL,
  body TEXT NOT NULL,
  score INTEGER NOT NULL DEFAULT 0,
  comments INTEGER NOT NULL DEFAULT 0,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE user_post_upvote(
  id INTEGER PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES user(id) ON DELETE CASCADE,
  post_id INTEGER NOT NULL REFERENCES user_post(id) ON DELETE CASCADE,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  UNIQUE(user_id, post_id)
);

CREATE TABLE user_post_downvote(
  id INTEGER PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES user(id) ON DELETE CASCADE,
  post_id INTEGER NOT NULL REFERENCES user_post(id) ON DELETE CASCADE,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  UNIQUE(user_id, post_id)
);

CREATE TABLE user_post_comment(
  id INTEGER PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES user(id) ON DELETE CASCADE,
  post_id INTEGER NOT NULL REFERENCES user_post(id) ON DELETE CASCADE,
  parent INTEGER REFERENCES user_post_comment(id) ON DELETE CASCADE,
  text TEXT NOT NULL,
  score INTEGER NOT NULL DEFAULT 0,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE user_post_comment_upvote(
  id INTEGER PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES user(id) ON DELETE CASCADE,
  post_comment_id INTEGER NOT NULL REFERENCES user_post(id) ON DELETE CASCADE,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  UNIQUE(user_id, post_comment_id)
);

CREATE TABLE user_post_comment_downvote(
  id INTEGER PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES user(id) ON DELETE CASCADE,
  post_comment_id INTEGER NOT NULL REFERENCES user_post(id) ON DELETE CASCADE,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  UNIQUE(user_id, post_comment_id)
);

-- Trigger to insert a 'user_post_upvote' on behalf of the user that created a 'user_post'.

CREATE TRIGGER user_post_insert AFTER INSERT ON user_post
  BEGIN
    INSERT INTO user_post_upvote (user_id, post_id) VALUES (NEW.user_id, NEW.id);
  END;

-- Triggers for inserts and deletes on 'user_post_upvote' and 'user_post_downvote'.

CREATE TRIGGER user_post_upvote_insert AFTER INSERT ON user_post_upvote FOR EACH ROW
  BEGIN
    UPDATE user_post SET score = score + 1 WHERE id = NEW.post_id;
    DELETE FROM user_post_downvote WHERE post_id = NEW.post_id;
  END;

CREATE TRIGGER user_post_upvote_delete DELETE ON user_post_upvote
  BEGIN
    UPDATE user_post SET score = score - 1 WHERE id = OLD.post_id;
  END;

CREATE TRIGGER user_post_downvote_insert INSERT ON user_post_downvote
  BEGIN
    UPDATE user_post SET score = score - 1 WHERE id = NEW.post_id;
    DELETE FROM user_post_upvote WHERE user_id = NEW.user_id AND post_id = NEW.post_id;
  END;

CREATE TRIGGER user_post_downvote_delete DELETE ON user_post_downvote
  BEGIN
    UPDATE user_post SET score = score + 1 WHERE id = OLD.post_id;
  END;

-- Trigger to insert a 'user_post_comment_upvote' on behalf of the user that created a 'user_post_comment'.

CREATE TRIGGER user_post_comment_insert AFTER INSERT ON user_post_comment
  BEGIN
    INSERT INTO user_post_comment_upvote (user_id, post_comment_id) VALUES (NEW.user_id, NEW.id);
    UPDATE user_post SET comments = comments + 1 WHERE id = NEW.post_id;
  END;

-- Triggers for inserts and deletes on 'user_post_comment_upvote' and 'user_post_comment_downvote'.

CREATE TRIGGER user_post_comment_upvote_insert INSERT ON user_post_comment_upvote
  BEGIN
    UPDATE user_post_comment SET score = score + 1 WHERE id = NEW.post_comment_id;
    DELETE FROM user_post_comment_downvote WHERE user_id = NEW.user_id AND post_id = NEW.post_id;
  END;

CREATE TRIGGER user_post_comment_upvote_delete DELETE ON user_post_comment_upvote
  BEGIN
    UPDATE user_post_comment SET score = score - 1 WHERE id = OLD.post_comment_id;
  END;

CREATE TRIGGER user_post_comment_downvote_insert INSERT ON user_post_comment_downvote
  BEGIN
    UPDATE user_post_comment SET score = score - 1 WHERE id = NEW.post_comment_id;
    DELETE FROM user_post_comment_upvote WHERE user_id = NEW.user_id AND post_id = NEW.post_id;
  END;

CREATE TRIGGER user_post_comment_downvote_delete DELETE ON user_post_comment_downvote
  BEGIN
    UPDATE user_post_comment SET score = score + 1 WHERE id = OLD.post_comment_id;
  END;
