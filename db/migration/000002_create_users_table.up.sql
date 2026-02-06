CREATE TABLE IF NOT EXISTS public.users (
    account_id integer NOT NULL,
    name character varying(255),
    avatar character varying(255),
    bio text,
    deleted_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_name_key UNIQUE (name);

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_pkey PRIMARY KEY (account_id);

ALTER TABLE ONLY public.users
    ADD CONSTRAINT users_account_id_fkey FOREIGN KEY (account_id) REFERENCES public.accounts(id);
