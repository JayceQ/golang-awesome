# video_server

create table comments (
    id varchar(64) not null,
    video_id varchar(64) ,
    content text,
    time datetime default current_timestamp ,
    primary key (id)
);

create table sessions (
    session_id tinytext not null,
    TTL tinytext,
    login_name text
);

alter table sessions add primary key (session_id(64));

create table users (
    id int unsigned not null auto_increment,
    login_name varchar(64),
    pwd text not null,
    unique key(login_name),
    primary key (id)
);

create table video_del_rec (
    video_id varchar(64) not null,
    primary key (video_id)
);

create table video_info (
    id varchar(64) not null,
    author_id int(10),
    name text,
    display_ctime text,
    create_time datetime default current_timestamp ,
    primary key (id)
);