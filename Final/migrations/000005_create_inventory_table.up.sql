CREATE TABLE IF NOT EXISTS inventory (
                                      ID integer PRIMARY KEY,
                                      user_id bigint NOT NULL REFERENCES users ON DELETE CASCADE,
                                      items JSONB NOT NULL
                                                              );