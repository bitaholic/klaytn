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
  DEFAULT CHARSET = ascii