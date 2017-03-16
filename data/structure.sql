use dkvgo;

CREATE TABLE IF NOT EXISTS `user` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `username` varchar(20) NOT NULL DEFAULT '',
  `email` varchar(30) NOT NULL DEFAULT '',
  `password` varchar(50) NOT NULL DEFAULT '',
  `last_login_ip` varchar(30) NOT NULL DEFAULT '',
  `last_login_time` datetime DEFAULT NULL,
  `create_at` datetime NOT NULL,
  `update_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `username` (`username`),
  UNIQUE KEY `email` (`email`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `job` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `name` varchar(50) NOT NULL DEFAULT '',
  `video_dir` varchar(512) NOT NULL DEFAULT '',
  `output_dir` varchar(512) NOT NULL DEFAULT '',
  `start_frame` int(11) NOT NULL DEFAULT '0',
  `end_frame` int(11) NOT NULL DEFAULT '0',
  `algorithm` varchar(20) NOT NULL DEFAULT '',
  `priority` int(11) NOT NULL DEFAULT '0',
  `camera_type` varchar(10) NOT NULL DEFAULT '',
  `quality` varchar(10) NOT NULL DEFAULT '',
  `enable_top` varchar(1) NOT NULL DEFAULT '1',
  `enable_bottom` varchar(1) NOT NULL DEFAULT '1',
  `enable_color_adjust` varchar(1) NOT NULL DEFAULT '1',
  `save_debug_img` varchar(5) NOT NULL DEFAULT 'false',
  `status` int(11) NOT NULL DEFAULT '0',
  `progress` double NOT NULL DEFAULT '0',
  `creator_id` int(11) NOT NULL,
  `operator_id` int(11) NOT NULL,
  `create_at` datetime NOT NULL,
  `update_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  KEY `job_creator_id` (`creator_id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS `job_state` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `job_id` int(11) NOT NULL,
  `content` longtext NOT NULL,
  `create_at` datetime NOT NULL,
  `update_at` datetime NOT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `job_id` (`job_id`)
) ENGINE=InnoDB AUTO_INCREMENT=2 DEFAULT CHARSET=utf8;

insert into user values(null, 'admin', 'admin@visiondk.com', '32cd9fbfcde4a88327a6afdd5efd9f82', '', null, '2016-03-12 20:39:42', '2016-03-12 20:39:42');
-- insert into job values (null, 'test', '/data/video_dir/test', '/data/output_dir/test', '1200', '1250', '3D_AURA', '100', 'AURA', '4k', '1', '1', '1', 'true', '0', '0.0', '0', '0', '2016-03-09 16:51:42', '2016-03-09 16:51:42');

-- update job set status='0' where status='1' or status='2' 