CREATE TABLE `user`
(
    `id`              int          NOT NULL AUTO_INCREMENT COMMENT '用户id',
    `app`             varchar(32)  NOT NULL DEFAULT '' COMMENT '所属app',
    `username`        varchar(32)  NOT NULL DEFAULT '' COMMENT '用户名',
    `login_type`      varchar(16)  NOT NULL DEFAULT '' COMMENT '登录类型',
    `identify`        varchar(64)  NOT NULL DEFAULT '' COMMENT '标志性账号, 登录类型是email则为邮箱， mobile则为手机号',
    `password`        char(40)     NOT NULL DEFAULT '' COMMENT '密码',
    `nickname`        varchar(32)  NOT NULL DEFAULT '' COMMENT '昵称',
    `avatar`          varchar(255) NOT NULL DEFAULT '' COMMENT '头像',
    `create_time`     timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time`     timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    `last_login_time` timestamp    NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '最后登录时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE `uk_username` (`username`, `app`),
    UNIQUE `uk_identify` (`identify`, `app`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 100000
  DEFAULT CHARSET = utf8mb4 COMMENT ='用户表';


CREATE TABLE `oauth_user`
(
    `id`           int         NOT NULL AUTO_INCREMENT COMMENT 'id',
    `app`          varchar(32) NOT NULL DEFAULT '' COMMENT '所属app',
    `user_id`      varchar(32) NOT NULL DEFAULT '' COMMENT '第三方用户id',
    `login_type`   varchar(16) NOT NULL DEFAULT '' COMMENT '登录类型',
    `bind_user_id` int         NOT NULL DEFAULT 0 COMMENT '绑定的内部用户id(user.id)',
    `create_time`  timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `update_time`  timestamp   NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`) USING BTREE,
    UNIQUE `uk_identify` (`user_id`, `app`, `login_type`),
    INDEX `idx_bind_user` (`bind_user_id`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 100000
  DEFAULT CHARSET = utf8mb4 COMMENT ='Oauth用户表';

