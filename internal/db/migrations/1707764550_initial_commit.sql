CREATE TABLE IF NOT EXISTS users
(
    id uuid NOT NULL,
    first_name character varying(30),
    last_name character varying(30),
    email text,
    phone character varying(15),
    password text,
    username character varying(60),
    member_since date,
    next_renewal date,
    PRIMARY KEY (id),
    CONSTRAINT username UNIQUE (username)
);

CREATE TABLE IF NOT EXISTS mirroring_links
(
    id uuid NOT NULL,
    created_by_id uuid NOT NULL,
    nickname character varying(60),
    upload_date date,
    duration_ms bigint,
    PRIMARY KEY (id),
    CONSTRAINT created_by_id FOREIGN KEY (created_by_id)
        REFERENCES public.users (id) MATCH SIMPLE
        ON UPDATE CASCADE
        ON DELETE NO ACTION
        NOT VALID
);

CREATE TABLE IF NOT EXISTS files
(
    id uuid NOT NULL,
    name text NOT NULL,
    size_bytes bigint NOT NULL,
    upload_date date NOT NULL,
    mirror_link_id uuid,
    PRIMARY KEY (id),
    CONSTRAINT mirror_link_id FOREIGN KEY (mirror_link_id)
        REFERENCES public.mirroring_links (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE SET NULL
        NOT VALID
);

CREATE TABLE IF NOT EXISTS host_links
(
    mirror_id uuid NOT NULL,
    bunkr text,
    gofile text,
    pixeldrain text,
    cyberfile text,
    saint_to text,
    cyberdrop text,
    CONSTRAINT mirror_id FOREIGN KEY (mirror_id)
        REFERENCES public.mirroring_links (id) MATCH SIMPLE
        ON UPDATE NO ACTION
        ON DELETE NO ACTION
        NOT VALID
);