PRAGMA user_version = 1;

CREATE TABLE user(
  id INTEGER PRIMARY KEY,
  username TEXT UNIQUE NOT NULL,
  hashed_password BLOB NOT NULL,
  session_token TEXT,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE post(
  id INTEGER PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES user(id) ON DELETE CASCADE,
  title TEXT NOT NULL,
  url TEXT NOT NULL,
  body TEXT NOT NULL,
  score INTEGER NOT NULL DEFAULT 0,
  comments INTEGER NOT NULL DEFAULT 0,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE post_upvote(
  id INTEGER PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES user(id) ON DELETE CASCADE,
  post_id INTEGER NOT NULL REFERENCES post(id) ON DELETE CASCADE,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  UNIQUE(user_id, post_id)
);

CREATE TABLE post_downvote(
  id INTEGER PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES user(id) ON DELETE CASCADE,
  post_id INTEGER NOT NULL REFERENCES post(id) ON DELETE CASCADE,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  UNIQUE(user_id, post_id)
);

CREATE TABLE comment(
  id INTEGER PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES user(id) ON DELETE CASCADE,
  post_id INTEGER NOT NULL REFERENCES post(id) ON DELETE CASCADE,
  parent INTEGER REFERENCES comment(id) ON DELETE CASCADE,
  text TEXT NOT NULL,
  score INTEGER NOT NULL DEFAULT 0,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL
);

CREATE TABLE comment_upvote(
  id INTEGER PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES user(id) ON DELETE CASCADE,
  comment_id INTEGER NOT NULL REFERENCES post(id) ON DELETE CASCADE,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  UNIQUE(user_id, comment_id)
);

CREATE TABLE comment_downvote(
  id INTEGER PRIMARY KEY,
  user_id INTEGER NOT NULL REFERENCES user(id) ON DELETE CASCADE,
  comment_id INTEGER NOT NULL REFERENCES post(id) ON DELETE CASCADE,
  created_at DATETIME DEFAULT CURRENT_TIMESTAMP NOT NULL,
  UNIQUE(user_id, comment_id)
);

-- Trigger to insert a 'post_upvote' on behalf of the user that created a 'post'.

CREATE TRIGGER post_insert AFTER INSERT ON post
  BEGIN
    INSERT INTO post_upvote (user_id, post_id) VALUES (NEW.user_id, NEW.id);
  END;

-- Triggers for inserts and deletes on 'post_upvote' and 'post_downvote'.

CREATE TRIGGER post_upvote_insert AFTER INSERT ON post_upvote FOR EACH ROW
  BEGIN
    UPDATE post SET score = score + 1 WHERE id = NEW.post_id;
    DELETE FROM post_downvote WHERE post_id = NEW.post_id;
  END;

CREATE TRIGGER post_upvote_delete DELETE ON post_upvote
  BEGIN
    UPDATE post SET score = score - 1 WHERE id = OLD.post_id;
  END;

CREATE TRIGGER post_downvote_insert INSERT ON post_downvote
  BEGIN
    UPDATE post SET score = score - 1 WHERE id = NEW.post_id;
    DELETE FROM post_upvote WHERE user_id = NEW.user_id AND post_id = NEW.post_id;
  END;

CREATE TRIGGER post_downvote_delete DELETE ON post_downvote
  BEGIN
    UPDATE post SET score = score + 1 WHERE id = OLD.post_id;
  END;

-- Trigger to insert a 'comment_upvote' on behalf of the user that created a 'comment'.

CREATE TRIGGER comment_insert AFTER INSERT ON comment
  BEGIN
    INSERT INTO comment_upvote (user_id, comment_id) VALUES (NEW.user_id, NEW.id);
    UPDATE post SET comments = comments + 1 WHERE id = NEW.post_id;
  END;

-- Triggers for inserts and deletes on 'comment_upvote' and 'comment_downvote'.

CREATE TRIGGER comment_upvote_insert INSERT ON comment_upvote
  BEGIN
    UPDATE comment SET score = score + 1 WHERE id = NEW.comment_id;
    DELETE FROM comment_downvote WHERE user_id = NEW.user_id AND post_id = NEW.post_id;
  END;

CREATE TRIGGER comment_upvote_delete DELETE ON comment_upvote
  BEGIN
    UPDATE comment SET score = score - 1 WHERE id = OLD.comment_id;
  END;

CREATE TRIGGER comment_downvote_insert INSERT ON comment_downvote
  BEGIN
    UPDATE comment SET score = score - 1 WHERE id = NEW.comment_id;
    DELETE FROM comment_upvote WHERE user_id = NEW.user_id AND post_id = NEW.post_id;
  END;

CREATE TRIGGER comment_downvote_delete DELETE ON comment_downvote
  BEGIN
    UPDATE comment SET score = score + 1 WHERE id = OLD.comment_id;
  END;
