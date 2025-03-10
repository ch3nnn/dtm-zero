create database if not exists dtm_zero
    /*!40100 DEFAULT CHARACTER SET utf8mb4 */
;
-- ----------------------------
-- Table structure for order
-- ----------------------------
drop table if exists dtm_zero.order;
create table if not exists dtm_zero.order
(
    `id`        bigint NOT NULL AUTO_INCREMENT,
    `user_id`   bigint NOT NULL DEFAULT '0',
    `goods_id`  bigint NOT NULL DEFAULT '0' COMMENT '商品id',
    `num`       int    NOT NULL DEFAULT '0' COMMENT '下单数量',
    `row_state` int    NOT NULL DEFAULT '0' COMMENT '-1:下单回滚废弃 0:待支付',
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;


-- ----------------------------
-- Table structure for stock
-- ----------------------------
drop table if exists dtm_zero.stock;
create table if not exists dtm_zero.stock
(
    `id`       bigint NOT NULL AUTO_INCREMENT,
    `goods_id` bigint NOT NULL DEFAULT '0' COMMENT '商品id',
    `num`      int    NOT NULL DEFAULT '0' COMMENT '库存数量',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_goodsId` (`goods_id`)
) ENGINE = InnoDB
  AUTO_INCREMENT = 2
  DEFAULT CHARSET = utf8mb4
  COLLATE = utf8mb4_general_ci;

-- ----------------------------
-- Records of stock
-- ----------------------------
BEGIN;
INSERT INTO dtm_zero.stock
VALUES (1, 1, 100);
COMMIT;


-- ----------------------------
-- dtm_barrier
-- ----------------------------

create database if not exists dtm_barrier
    /*!40100 DEFAULT CHARACTER SET utf8mb4 */
;
drop table if exists dtm_barrier.barrier;
create table if not exists dtm_barrier.barrier
(
    id          bigint(22) PRIMARY KEY AUTO_INCREMENT,
    trans_type  varchar(45)  default '',
    gid         varchar(128) default '',
    branch_id   varchar(128) default '',
    op          varchar(45)  default '',
    barrier_id  varchar(45)  default '',
    reason      varchar(45)  default '' comment 'the branch type who insert this record',
    create_time datetime     DEFAULT now(),
    update_time datetime     DEFAULT now(),
    key (create_time),
    key (update_time),
    UNIQUE key (gid, branch_id, op, barrier_id)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;



-- ----------------------------
-- dtm
-- ----------------------------

CREATE DATABASE IF NOT EXISTS dtm
    /*!40100 DEFAULT CHARACTER SET utf8mb4 */
;
drop table IF EXISTS dtm.trans_global;
CREATE TABLE if not EXISTS dtm.trans_global
(
    `id`                 bigint(22)    NOT NULL AUTO_INCREMENT,
    `gid`                varchar(128)  NOT NULL COMMENT 'global transaction id',
    `trans_type`         varchar(45)   not null COMMENT 'transaction type: saga | xa | tcc | msg',
    `status`             varchar(12)   NOT NULL COMMENT 'transaction status: prepared | submitted | aborting | succeed | failed',
    `query_prepared`     varchar(1024) NOT NULL COMMENT 'url to check for msg|workflow',
    `protocol`           varchar(45)   not null comment 'protocol: http | grpc | json-rpc',
    `create_time`        datetime               DEFAULT NULL,
    `update_time`        datetime               DEFAULT NULL,
    `finish_time`        datetime               DEFAULT NULL,
    `rollback_time`      datetime               DEFAULT NULL,
    `options`            varchar(1024)          DEFAULT '' COMMENT 'options for transaction like: TimeoutToFail, RequestTimeout',
    `custom_data`        varchar(1024)          DEFAULT '' COMMENT 'custom data for transaction',
    `next_cron_interval` int(11)                default null comment 'next cron interval. for use of cron job',
    `next_cron_time`     datetime               default null comment 'next time to process this trans. for use of cron job',
    `owner`              varchar(128)  not null default '' comment 'who is locking this trans',
    `ext_data`           TEXT comment 'extra data for this trans. currently used in workflow pattern',
    `result`             varchar(1024)          DEFAULT '' COMMENT 'result for transaction',
    `rollback_reason`    varchar(1024)          DEFAULT '' COMMENT 'rollback reason for transaction',
    PRIMARY KEY (`id`),
    UNIQUE KEY `gid` (`gid`),
    key `owner` (`owner`),
    key `status_next_cron_time` (`status`, `next_cron_time`) comment 'cron job will use this index to query trans'
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;
drop table IF EXISTS dtm.trans_branch_op;
CREATE TABLE IF NOT EXISTS dtm.trans_branch_op
(
    `id`            bigint(22)    NOT NULL AUTO_INCREMENT,
    `gid`           varchar(128)  NOT NULL COMMENT 'global transaction id',
    `url`           varchar(1024) NOT NULL COMMENT 'the url of this op',
    `data`          TEXT COMMENT 'request body, depreceated',
    `bin_data`      BLOB COMMENT 'request body',
    `branch_id`     VARCHAR(128)  NOT NULL COMMENT 'transaction branch ID',
    `op`            varchar(45)   NOT NULL COMMENT 'transaction operation type like: action | compensate | try | confirm | cancel',
    `status`        varchar(45)   NOT NULL COMMENT 'transaction op status: prepared | succeed | failed',
    `finish_time`   datetime DEFAULT NULL,
    `rollback_time` datetime DEFAULT NULL,
    `create_time`   datetime DEFAULT NULL,
    `update_time`   datetime DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE KEY `gid_uniq` (`gid`, `branch_id`, `op`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;
drop table IF EXISTS dtm.kv;
CREATE TABLE IF NOT EXISTS dtm.kv
(
    `id`        bigint(22)   NOT NULL AUTO_INCREMENT,
    `cat`       varchar(45)  NOT NULL COMMENT 'the category of this data',
    `k`         varchar(128) NOT NULL,
    `v`         TEXT,
    `version`   bigint(22) default 1 COMMENT 'version of the value',
    create_time datetime   default NULL,
    update_time datetime   DEFAULT NULL,
    PRIMARY KEY (`id`),
    UNIQUE key `uniq_k` (`cat`, `k`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8mb4;
