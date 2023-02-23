CREATE TABLE IF NOT EXISTS cases (
                                     id bigserial PRIMARY KEY,
                                     name text NOT NULL,
                                     price integer NOT NULL,
                                     items JSONB NOT NULL
);