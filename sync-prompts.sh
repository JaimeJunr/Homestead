#!/bin/bash

# ==============================================================================
# Backup / restore de .claude entre o projeto atual e o repositório central (prompts-ai)
#
#   --backup   Projeto → central
#   --restore  Central → projeto
# ==============================================================================

set -euo pipefail

readonly SCRIPT_NAME="$(basename "$0")"
readonly CLAUDE_DIR=".agents"

# Resolução do diretório central (absoluto):
#   1) SYNC_PROMPTS_CENTRAL
#   2) diretório do script, se for o repo central (contém .claude)
#   3) ../prompts-ai a partir do cwd
resolve_central_base() {
  local script_dir
  script_dir="$(cd "$(dirname "$0")" && pwd)"
  if [[ -n "${SYNC_PROMPTS_CENTRAL:-}" ]]; then
    if [[ ! -d "${SYNC_PROMPTS_CENTRAL}" ]]; then
      echo "❌ Erro: SYNC_PROMPTS_CENTRAL não é um diretório: ${SYNC_PROMPTS_CENTRAL}" >&2
      exit 1
    fi
    if [[ ! -d "${SYNC_PROMPTS_CENTRAL}/${CLAUDE_DIR}" ]]; then
      echo "❌ Erro: SYNC_PROMPTS_CENTRAL não parece o repo prompts-ai (falta ${CLAUDE_DIR}): ${SYNC_PROMPTS_CENTRAL}" >&2
      exit 1
    fi
    echo "$(cd "$SYNC_PROMPTS_CENTRAL" && pwd)"
    return
  fi
  if [[ -d "${script_dir}/${CLAUDE_DIR}" ]]; then
    echo "$script_dir"
    return
  fi
  echo "$(cd "$(pwd)/../prompts-ai" 2>/dev/null && pwd)" || echo ""
}

CENTRAL_BASE="$(resolve_central_base)"
if [[ -z "$CENTRAL_BASE" || ! -d "$CENTRAL_BASE" ]]; then
  echo "❌ Erro: Diretório central não encontrado. Defina SYNC_PROMPTS_CENTRAL ou execute a partir do repo prompts-ai ou de um projeto com ../prompts-ai." >&2
  exit 1
fi
readonly CENTRAL_BASE

readonly PROJECT_CLAUDE_PATH="./${CLAUDE_DIR}/"
readonly CENTRAL_CLAUDE_PATH="${CENTRAL_BASE}/${CLAUDE_DIR}/"

error_exit() {
  echo "❌ Erro: $1" >&2
  exit 1
}

success() { echo "✅ $1"; }
info() { echo "🔄 $1"; }
process() { echo "📁 $1"; }

show_help() {
  cat << EOF
Uso: $SCRIPT_NAME --backup|--restore [opções]

Modos:
  --backup    Copia ${CLAUDE_DIR}/ do projeto atual para o repositório central
  --restore   Copia ${CLAUDE_DIR}/ do repositório central para o projeto atual

Opções:
  --delete    Remove no destino o que não existe na origem (espelhar)
  --dry-run   Simula sem alterar arquivos
  --help      Esta ajuda

Variável de ambiente:
  SYNC_PROMPTS_CENTRAL   Path absoluto do clone do prompts-ai (opcional)

Exemplos:
  $SCRIPT_NAME --backup
  $SCRIPT_NAME --restore
  $SCRIPT_NAME --restore --delete
  $SCRIPT_NAME --restore --dry-run
EOF
}

parse_arguments() {
  local mode=""
  local rsync_extra=()
  while [[ $# -gt 0 ]]; do
    case "$1" in
      --backup)  mode="backup"; shift ;;
      --restore) mode="restore"; shift ;;
      --delete)  rsync_extra+=("--delete"); shift ;;
      --dry-run) rsync_extra+=("--dry-run"); shift ;;
      --help)    show_help; exit 0 ;;
      *)         error_exit "Opção desconhecida '$1'. Use --help." ;;
    esac
  done
  [[ -z "$mode" ]] && error_exit "Informe --backup ou --restore."
  MODE="$mode"
  RSYNC_EXTRA=("${rsync_extra[@]}")
}

execute_sync() {
  local origem destino
  if [[ "$MODE" == "backup" ]]; then
    origem="$PROJECT_CLAUDE_PATH"
    destino="$CENTRAL_CLAUDE_PATH"
    info "BACKUP: projeto → central"
  else
    origem="$CENTRAL_CLAUDE_PATH"
    destino="$PROJECT_CLAUDE_PATH"
    info "RESTORE: central → projeto"
  fi

  echo "    Central: $CENTRAL_BASE"
  echo "    Origem:  $origem"
  echo "    Destino: $destino"
  echo ""

  if [[ ! -d "$origem" ]]; then
    if [[ "$MODE" == "restore" ]]; then
      error_exit "Pasta ${CLAUDE_DIR}/ não encontrada no repositório central."
    fi
    error_exit "Pasta ${CLAUDE_DIR}/ não encontrada no projeto atual."
  fi

  mkdir -p "$destino"
  process "rsync ${CLAUDE_DIR}/..."
  # Trailing slash: conteúdo de .claude → destino
  if rsync -av "${RSYNC_EXTRA[@]}" "$origem" "$destino"; then
    if [[ " ${RSYNC_EXTRA[*]} " == *" --dry-run "* ]]; then
      success "Simulação (--dry-run) concluída."
    else
      success "Sincronização de ${CLAUDE_DIR}/ concluída."
    fi
  else
    error_exit "Falha no rsync."
  fi
}

main() {
  parse_arguments "$@"
  execute_sync
}

main "$@"
