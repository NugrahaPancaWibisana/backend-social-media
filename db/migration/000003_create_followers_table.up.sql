CREATE TABLE IF NOT EXISTS public.followers (
    following_user_id integer,
    followed_user_id integer,
    deleted_at timestamp without time zone,
    created_at timestamp without time zone DEFAULT now() NOT NULL,
    updated_at timestamp without time zone DEFAULT now() NOT NULL
);

ALTER TABLE ONLY public.followers
    ADD CONSTRAINT followers_followed_user_id_fkey FOREIGN KEY (followed_user_id) REFERENCES public.users(account_id);

ALTER TABLE ONLY public.followers
    ADD CONSTRAINT followers_following_user_id_fkey FOREIGN KEY (following_user_id) REFERENCES public.users(account_id);
