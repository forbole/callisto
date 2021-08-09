#!/bin/bash
set -e

bdjuno_init () {
  flags=""
  if [[ ! -z "${CLIENT_NAME}" ]]; then
   flags=" ${flags} --client-name ${CLIENT_NAME}"
  fi
  if [[ ! -z "${COSMOS_MODULE}" ]]; then
   flags=" ${flags} --cosmos-modules ${COSMOS_MODULE}"
  fi
  if [[ ! -z "${COSMOS_PREFIX}" ]]; then
   flags=" ${flags} --cosmos-prefix ${COSMOS_PREFIX}"
  fi
  if [[ ! -z "${DATABASE_NAME}" ]]; then
   flags=" ${flags} --database-name ${DATABASE_NAME}"
  fi
  if [[ ! -z "${DATABASE_HOST}" ]]; then
   flags=" ${flags} --database-host ${DATABASE_HOST}"
  fi
  if [[ ! -z "${DATABASE_PASSWORD}" ]]; then
   flags=" ${flags} --database-password ${DATABASE_PASSWORD}"
  fi
  if [[ ! -z "${DATABASE_PORT}" ]]; then
   flags=" ${flags} --database-port ${DATABASE_PORT}"
  fi
  if [[ ! -z "${DATABASE_SCHEMA}" ]]; then
   flags=" ${flags} --database-schema ${DATABASE_SCHEMA}"
  fi
  if [[ ! -z "${DATABASE_SSL_MODE}" ]]; then
   flags=" ${flags} --database-ssl-mode ${DATABASE_SSL_MODE}"
  fi
  if [[ ! -z "${DATABASE_USER}" ]]; then
   flags=" ${flags} --database-user ${DATABASE_USER}"
  fi
  if [[ ! -z "${GRPC_ADDRESS}" ]]; then
   flags=" ${flags} --grpc-address ${GRPC_ADDRESS}"
  fi
  if [[ ! -z "${GRPC_INSECURE}" ]]; then
   flags=" ${flags} --grpc-insecure ${GRPC_INSECURE}"
  fi
  if [[ ! -z "${LOGGING_FORMAT}" ]]; then
   flags=" ${flags} --logging-format ${LOGGING_FORMAT}"
  fi
    if [[ ! -z "${LOGGING_LEVEL}" ]]; then
   flags=" ${flags} --logging-level ${LOGGING_LEVEL}"
  fi
    if [[ ! -z "${MAX_IDLE_CONNECTIONS}" ]]; then
   flags=" ${flags} --max-idle-connections ${MAX_IDLE_CONNECTIONS}"
  fi
    if [[ ! -z "${MAX_OPEN_CONNECTIONS}" ]]; then
   flags=" ${flags} --max-open-connections  ${MAX_OPEN_CONNECTIONS}"
  fi
    if [[ ! -z "${PARSING_FAST_SYNC}" ]]; then
   flags=" ${flags} --parsing-fast-sync  ${PARSING_FAST_SYNC}"
  fi
    if [[ ! -z "${PARSING_GENESIS_FILE_PATH}" ]]; then
   flags=" ${flags} --parsing-genesis-file-path ${PARSING_GENESIS_FILE_PATH}"
  fi
    if [[ ! -z "${PARSING_NEW_BLOCKS}" ]]; then
   flags=" ${flags} --parsing-new-blocks  ${PARSING_NEW_BLOCKS}"
  fi
  if [[ ! -z "${PARSING_OLD_BLOCKS}" ]]; then
   flags=" ${flags} --parsing-old-blocks ${PARSING_OLD_BLOCKS}"
  fi
  if [[ ! -z "${PARSING__PARSE_GENESIS}" ]]; then
   flags=" ${flags} --parsing-parse-genesis ${PARSING__PARSE_GENESIS}"
  fi
  if [[ ! -z "${PARSING_START_HEIGHT}" ]]; then
   flags=" ${flags} --parsing-start-height ${PARSING_START_HEIGHT}"
  fi
  if [[ ! -z "${PARSING_WORKERS}" ]]; then
   flags=" ${flags} --parsing-workers ${PARSING_WORKERS}"
  fi
  if [[ ! -z "${PRUNING_INTERVAL}" ]]; then
   flags=" ${flags} --pruning-interval ${PRUNING_INTERVAL}"
  fi
  if [[ ! -z "${PRUNING_KEEP_EVERY}" ]]; then
   flags=" ${flags} --pruning-keep-every ${PRUNING_KEEP_EVERY}"
  fi
  if [[ ! -z "${PRUNING_KEEP_RECENT}" ]]; then
   flags=" ${flags} --pruning-keep-recent ${PRUNING_KEEP_RECENT}"
  fi
  if [[ ! -z "${RPC_ADDRESS}" ]]; then
   flags=" ${flags} --rpc-address ${RPC_ADDRESS}"
  fi
  if [[ ! -z "${TELEMETRY_ENABLED}" ]]; then
   flags=" ${flags} --telemetry-enabled ${TELEMETRY_ENABLED}"
  fi
  if [[ ! -z "${TELEMETRY_PORT}" ]]; then
   flags=" ${flags} --telemetry-port ${TELEMETRY_PORT}"
  fi
  
  bdjuno init --home ${BDJUNO_HOME} ${flags}
}
# Default CMD
if [ "$1" = 'bdjuno' ] && [[ -z "$2" ]] ; then
  BDJUNO_HOME=${BDJUNO_HOME:-/bdjuno/.bdjuno}
  bdjuno_init
  bdjuno parse --home ${BDJUNO_HOME}
else
  # This allow user to use other commands
  exec "$@"
fi
