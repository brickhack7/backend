create table profiles(
    uid text primary key, 
    -- bday text not null,
    lat real NOT NULL,
    long real NOT NULL,
    interests text[]
);

create table locations (
    loc_id text not null PRIMARY KEY,
    name text,
    city text,
    lat real NOT NULL,
    long real NOT NULL
);

create table user_matches (
    uid1 text not null,
    uid2 text not null,
    status int not null,
    PRIMARY KEY (uid1, uid2),
    CONSTRAINT not_self check (uid1 <> uid2),
    CONSTRAINT fk_uid1 FOREIGN KEY (uid1) REFERENCES profiles(uid),
    CONSTRAINT fk_uid2 FOREIGN KEY (uid2) REFERENCES profiles(uid)
    ON DELETE CASCADE
);

create table loc_matches (
    uid text not null,
    loc_id text not null,
    status int not null,
    PRIMARY KEY (uid, loc_id),
    CONSTRAINT fk_uid FOREIGN KEY (uid) REFERENCES profiles(uid),
    CONSTRAINT fk_locid FOREIGN KEY (loc_id) REFERENCES locations(loc_id)
    ON DELETE CASCADE
);

-- CREATE TABLE messages (
--     uuid UUID PRIMARY KEY,
--     msg_from text,
--     msg_to text,
--     message text,
--     time text,
--     CONSTRAINT fk_from_uid FOREIGN KEY (msg_from) REFERENCES profiles(uid),
--     CONSTRAINT fk_to_uid FOREIGN KEY (msg_to) REFERENCES profiles(uid)
-- );