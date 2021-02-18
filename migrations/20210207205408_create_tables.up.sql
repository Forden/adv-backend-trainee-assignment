create table if not exists ads
(
    ad_id       text not null
    constraint ads_pkey
    primary key,
    title       text,
    description text,
    price       integer,
    photo_links text,
    created_at  integer
);