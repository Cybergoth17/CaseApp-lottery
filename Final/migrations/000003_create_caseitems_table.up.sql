CREATE TABLE IF NOT EXISTS caseitems (
                                         id bigserial PRIMARY KEY,
                                         itemname text NOT NULL,
                                         itemdesc text NOT NULL,
                                         type text NOT NULL,
                                         stars integer NOT NULL
);