--user tab
CREATE TABLE `tbl_user`
(
    `id`          int(11) NOT NULL AUTO_INCREMENT,
    `username`    varchar(256) DEFAULT '' COMMENT '',
    `user_pwd` varchar(256)  DEFAULT '' COMMENT '',
    `openid`      varchar(256) DEFAULT '' COMMENT '',
    `unionid`     varchar(256) DEFAULT '' COMMENT '',
    `phone`       varchar(256) NOT NULL DEFAULT '' COMMENT '',
    `session_key` varchar(256) DEFAULT '' COMMENT '',
    `foreverid`   varchar(256) DEFAULT '' COMMENT '',
    `totalstor`      bigint(20) DEFAULT '0' COMMENT '',
    `file_num`       bigint(20) DEFAULT '0' COMMENT '',
    `starfile_num`       bigint(20) DEFAULT '0' COMMENT '',
    `creat_num`      bigint(20) DEFAULT '0' COMMENT '',
    `child_num`      bigint(20) DEFAULT '0' COMMENT '',
    `young_num`      bigint(20) DEFAULT '0' COMMENT '',
    `middle_num`      bigint(20) DEFAULT '0' COMMENT '',
    `old_num`        bigint(20) DEFAULT '0' COMMENT '',
    `nick_name`   varchar(256) DEFAULT '' COMMENT '',
    `avator_url`  varchar(256) DEFAULT '' COMMENT '',
    `city`        varchar(64)           DEFAULT '' COMMENT '',
    `province`    varchar(64)           DEFAULT '' COMMENT '',
    `country`     varchar(64)           DEFAULT '' COMMENT '',
    `gender`      int(11) DEFAULT '0' COMMENT '',
    `integral`      bigint(20) DEFAULT '0' COMMENT '',
    `status`      int(11) DEFAULT '0' COMMENT '',
    `created_at`  bigint(20)  DEFAULT '0' COMMENT '',
    `update_at`   bigint(20)  DEFAULT '0' COMMENT '',
    `profile`     text COMMENT '',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_phone` (`phone`),
    KEY           `idx_status` (`status`)
) ENGINE=InnoDB AUTO_INCREMENT=5 DEFAULT CHARSET=utf8;

-- user token tab
CREATE TABLE `tbl_user_token`
(
    `id`         int(11) NOT NULL AUTO_INCREMENT,
    `phone`   varchar(256) NOT NULL DEFAULT '' COMMENT '',
    `user_token` char(255)     NOT NULL DEFAULT '' COMMENT '',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_phone` (`phone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8


-- user file tab
CREATE TABLE `tbl_user_file`
(
    `id`         int(11) NOT NULL PRIMARY KEY AUTO_INCREMENT,
    `phone`    varchar(256) DEFAULT '' COMMENT '',
    `cid`            varchar(256)  NOT NULL DEFAULT '' COMMENT '',
    `file_name`  varchar(256) DEFAULT '' COMMENT '',
    `ext`  varchar(256) DEFAULT '' COMMENT '',
    `file_type`      varchar(256)  DEFAULT '' COMMENT '',
    `file_size`  bigint(20) DEFAULT '0' COMMENT '',
    `created_at`  bigint(20)  DEFAULT '0' COMMENT '',
    `comment`            varchar(256)  DEFAULT '' COMMENT '',
    `folder_type`      int(11) DEFAULT '0' COMMENT '',
    `star`  int(11) DEFAULT '0' COMMENT '',
    `tags`            varchar(256)  DEFAULT '' COMMENT '',
    `minio_url`            varchar(256)  DEFAULT '' COMMENT '',
    `status`     int(11) DEFAULT '0' COMMENT '',
    `update_at`  bigint(20)  DEFAULT '0' COMMENT ''
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE `tbl_wx_request`
(
    `id`         int(11) NOT NULL AUTO_INCREMENT,
    `openid`      varchar(256) NOT NULL DEFAULT '' COMMENT '',
    `unionid` varchar(256) DEFAULT '' COMMENT '',
    `session_key` varchar(256) NOT NULL DEFAULT '' COMMENT '',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_openid` (`openid`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

