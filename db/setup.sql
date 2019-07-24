CREATE TABLE `potentials` (
  `id`               BIGINT(20)    NOT NULL AUTO_INCREMENT,
  `date`             DATETIME      NOT NULL DEFAULT CURRENT_TIMESTAMP,

  `request_method`   VARCHAR(50)   NOT NULL,
  `request_url`      VARCHAR(2000) NOT NULL,
  `request_payload`  VARCHAR(2000) NOT NULL,
  `response_status`  INTEGER(10)   NOT NULL,
  `response_headers` VARCHAR(5000) NOT NULL,
  `response_payload` MEDIUMTEXT    NOT NULL,

  PRIMARY KEY (`id`, `date`),
  KEY `search_request` (`request_method`),
  KEY `search_response` (`response_status`)
)
  ENGINE = InnoDB
  DEFAULT CHARSET = utf8
  PARTITION BY LIST (month(date))
  (PARTITION p1 VALUES IN (1)
    ENGINE = InnoDB,
  PARTITION p2 VALUES IN (2)
    ENGINE = InnoDB,
  PARTITION p3 VALUES IN (3)
    ENGINE = InnoDB,
  PARTITION p4 VALUES IN (4)
    ENGINE = InnoDB,
  PARTITION p5 VALUES IN (5)
    ENGINE = InnoDB,
  PARTITION p6 VALUES IN (6)
    ENGINE = InnoDB,
  PARTITION p7 VALUES IN (7)
    ENGINE = InnoDB,
  PARTITION p8 VALUES IN (8)
    ENGINE = InnoDB,
  PARTITION p9 VALUES IN (9)
    ENGINE = InnoDB,
  PARTITION p10 VALUES IN (10)
    ENGINE = InnoDB,
  PARTITION p11 VALUES IN (11)
    ENGINE = InnoDB,
  PARTITION p12 VALUES IN (12)
    ENGINE = InnoDB);

SELECT *
FROM potentials;

SELECT
  request_method,
  request_url,
  response_status,
  count(*)
FROM potentials
GROUP BY request_method, request_url, response_status;

TRUNCATE TABLE potentials;