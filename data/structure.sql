use dkvgo;

create table if not exists user (
    id int not null auto_increment,
    username varchar(20) not null,
    email varchar(30) not null,
    password varchar(50) not null,
    create_at int not null,
    update_at int not null,
    primary key (id),
    unique email_idx (email)
    unique username_idx(username)
) engine=innodb;

create table  if not exists job (
    id int not null auto_increment,
    name varchar(50) not null,
    video_dir varchar(512) not null,
    output_dir varchar(512) not null,
    start_frame int not null,
    end_frame int not null,
    algorithm varchar(20) not null,
    priority int not null default 100,
    camera_type varchar(10) not null,
    enable_top char(10) not null default '1',
    enable_bottom char(10) not null default '1',
    enable_color_adjust char(1) not null default '1',
    quality varchar(10) not null,
    save_debug_img varchar(10) not null default 'false',
    status int not null default 0,
    progress float not null default 0.0,
    creator_id int not null default 0,
    operator_id int not null default 0,
    create_at datetime not null,
    update_at datetime not null,
    primary key (id)
) engine=innodb;

create table if not exists job_state (
    id int not null auto_increment,
    job_id int not null,
    content text not null,
    create_at datetime not null,
    update_at datetime not null,
    primary key (id),
    index job_id_idx (job_id)
) engine=innodb;

-- insert into user values();
insert into job values (null, 'test', '/data/video_dir/test', '/data/output_dir/test', '1200', '1250', '3D_AURA', '100', 'AURA', '1', '1', '1', '4k', 'true', '0', '0.0', '0', '0', '2016-03-09 16:51:42', '2016-03-09 16:51:42');

update job set status='0' where status='1' or status='2' 