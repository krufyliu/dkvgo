use dkvgo;

create table if not exists users (
    id int not null auto_increment,
    username varchar(20) not null,
    email varchar(30) not null,
    password varchar(50) not null,
    create_at int not null,
    update_at int not null,
    primary key (id),
    unique email_idx (email)
) engine=innodb;

create table  if not exists jobs (
    id int not null auto_increment,
    name varchar(50) not null,
    video_dir varchar(512) not null,
    output_dir varchar(512) not null,
    start_frame int not null,
    end_frame int not null,
    algorithm varchar(20) not null,
    priority int not null default 100,
    status int not null default 0,
    camera_type varchar(10) not null,
    enable_top char(1) not null,
    enable_bottom char(1) not null,
    enable_color_adjust char(1) not null default '1',
    quality varchar(10) not null,
    progress float not null default 0.0,
    create_at int not null,
    update_at int not null,
    creator_id int not null default 0,
    operator_id int not null default 0,
    primary key (id)
) engine=innodb;

create table if not exists job_states (
    id int not null auto_increment,
    job_id int not null,
    content text not null,
    create_at int not null,
    update_at int not null,
    primary key (id),
    index job_id_idx (job_id)
) engine=innodb;

insert into jobs values (null, 'test', '/data/video_dir/test', '/data/output_dir/test', '1200', '1250', 'FACEBOOK_3D', '100', '0', 'GOPRO', '1', '1', '1', '4k', '0.0', '1488968949', '1488968949');

update jobs set status='0' where status='1' or status='2' 