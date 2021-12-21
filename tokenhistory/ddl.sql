CREATE TABLE `token_history`
(
    `block_num`           int unsigned DEFAULT NULL,
    `tx_idx`              int unsigned DEFAULT NULL,
    `log_idx`             int unsigned DEFAULT NULL,
    `to_addr`             binary(20)   DEFAULT NULL,
    `from_addr`           binary(20)   DEFAULT NULL,
    `tx_value`            bigint       DEFAULT NULL,
    `token_contract_addr` binary(20)   DEFAULT NULL,
    KEY `token_history_from_addr_block_num_index` (`from_addr`, `block_num`),
    KEY `token_history_to_addr_block_num_index` (`to_addr`, `block_num`)
) ENGINE = InnoDB
  DEFAULT CHARSET = ascii;

CREATE TABLE `klay_transfer_history`
(
    `account_addr`  binary(20)   NOT NULL,
    `block_num`     int unsigned NOT NULL,
    `tx_idx`        int unsigned NOT NULL,
    `itx_idx`       int unsigned NOT NULL,
    `opposite_addr` binary(20)  DEFAULT NULL,
    `tx_value`      varchar(80) DEFAULT NULL,
    `balance`       varchar(80) DEFAULT NULL,
    `tx_hash`       binary(32)   NOT NULL,
    `direction`     tinyint(1)  DEFAULT NULL,
    `block_time`    int unsigned NOT NULL,
    KEY `token_history_from_addr_block_num_index` (`account_addr`, `block_num`, `tx_idx`, `itx_idx`),
    KEY `account_block_time_index` (`account_addr`, `block_time` DESC)
) ENGINE = InnoDB
  DEFAULT CHARSET = ascii;

