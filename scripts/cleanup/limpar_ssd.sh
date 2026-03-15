#!/bin/bash

# Script Orquestrador de Limpeza de SSD
# Aplica o Princípio da Responsabilidade Única (SOLID)

set -e

# Detectar usuário real
REAL_USER=${SUDO_USER:-$USER}
REAL_HOME=$(getent passwd "$REAL_USER" | cut -d: -f6)
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

echo "=== Limpeza de Espaço no SSD (Orquestrador) ==="
echo "Usuário: $REAL_USER"
echo "Home: $REAL_HOME"
echo ""

# --- Funções Utilitárias ---
show_space() {
    echo "Espaço disponível:"
    df -h / | tail -1
    echo ""
}

confirm_action() {
    local description="$1"
    local command_to_run="$2"
    local size_check_cmd="$3"

    echo "---------------------------------------------------"
    echo "$description"
    
    if [ -n "$size_check_cmd" ]; then
        local size=$(eval "$size_check_cmd" 2>/dev/null || echo "N/A")
        echo "   Tamanho estimado: $size"
    fi

    read -p "   Deseja prosseguir com esta limpeza? [s/N]: " -n 1 -r
    echo ""
    if [[ $REPLY =~ ^[Ss]$ ]]; then
        echo "   Executando..."
        eval "$command_to_run" || echo "   ⚠ Erro ao executar (pode precisar de sudo ou não existir)"
        echo "   ✓ Concluído."
    else
        echo "   ⊘ Pular."
    fi
    echo ""
}

# Exportar variáveis e funções para sub-scripts
export REAL_USER REAL_HOME
export -f confirm_action

# --- Execução ---
show_space

# 1. Limpeza Geral (Caches, Sistema, etc)
if [ -f "$SCRIPT_DIR/limpar_geral.sh" ]; then
    source "$SCRIPT_DIR/limpar_geral.sh"
else
    echo "⚠ Erro: limpar_geral.sh não encontrado."
fi

# 2. Limpeza de Itens Grandes (Arquivos e Pastas)
if [ -f "$SCRIPT_DIR/limpar_grandes.sh" ]; then
    source "$SCRIPT_DIR/limpar_grandes.sh"
else
    echo "⚠ Erro: limpar_grandes.sh não encontrado."
fi

show_space
echo "=== Limpeza Concluída! ==="
