#!/bin/bash

# Script para Limpar Pastas e Arquivos Grandes
# Parte do sistema de limpeza SOLID

# Fallback para execução independente
REAL_USER=${REAL_USER:-${SUDO_USER:-$USER}}
REAL_HOME=${REAL_HOME:-$(getent passwd "$REAL_USER" | cut -d: -f6)}

echo "--- Iniciando Busca por Itens Grandes ---"

# 15. Large Files Scanner
echo "---------------------------------------------------"
echo "15. Buscar arquivos grandes (>100MB) no Home?"
read -p "    Deseja iniciar a busca? [s/N]: " -n 1 -r
echo ""
if [[ $REPLY =~ ^[Ss]$ ]]; then
    echo "    Buscando arquivos maiores que 100MB (ignorando pastas ocultas e específicas)..."
    TMP_LIST=$(mktemp)
    find "$REAL_HOME" -maxdepth 4 -not -path '*/.*' -not -path '*node_modules*' -not -path '*databases/datamysql*' -type f -size +100M -printf "%s\t%p\n" 2>/dev/null | sort -rn | head -n 50 > "$TMP_LIST"
    
    TOTAL_FOUND=$(wc -l < "$TMP_LIST")
    if [ "$TOTAL_FOUND" -eq 0 ]; then
        echo "    Nenhum arquivo grande encontrado."
    else
        echo "    Encontrados $TOTAL_FOUND arquivos grandes (Top 50)."
        echo "    Opções: [s]im, [n]ão, [q]uit"
        echo ""
        
        while read size_bytes filepath; do
            if [ ! -f "$filepath" ]; then continue; fi
            
            human_size=$(du -h "$filepath" | cut -f1)
            
            echo "    Arquivo: $filepath"
            echo "    Tamanho: $human_size"
            
            read -p "    >> Deletar? [s/n/q]: " -n 1 -r < /dev/tty
            echo ""
            
            case $REPLY in
                [sS])
                    rm -v "$filepath"
                    echo "       ✓ Deletado."
                    ;;
                [qQ])
                    echo "       Parando..."
                    break
                    ;;
                *)
                    echo "       ⊘ Mantido."
                    ;;
            esac
            echo "---------------------------------------------------"
        done < "$TMP_LIST"
    fi
    rm -f "$TMP_LIST"
else
    echo "    Busca pulada."
fi


# 16. Large Directories / Games
echo "---------------------------------------------------"
echo "16. Buscar pastas grandes (>1GB) ou Jogos?"
read -p "    Deseja listar pastas gigantes? [s/N]: " -n 1 -r
echo ""
if [[ $REPLY =~ ^[Ss]$ ]]; then
    echo "    Calculando tamanhos (isso pode levar um tempo)..."
    
    TMP_DIRS=$(mktemp)
    touch "$TMP_DIRS"

    # 1. Steam Games
    if [ -d "$REAL_HOME/.local/share/Steam/steamapps/common" ]; then
         echo "    Verificando Steam..."
         du -h --max-depth=1 "$REAL_HOME/.local/share/Steam/steamapps/common" 2>/dev/null | sort -hr >> "$TMP_DIRS"
    fi
    
    if [ -d "$REAL_HOME/.steam/steam/steamapps/common" ]; then
         echo "    Verificando Steam Legacy..."
         du -h --max-depth=1 "$REAL_HOME/.steam/steam/steamapps/common" 2>/dev/null | sort -hr >> "$TMP_DIRS"
    fi

    # 3. User Games folder
    if [ -d "$REAL_HOME/Games" ]; then
         echo "    Verificando ~/Games..."
         du -h --max-depth=1 "$REAL_HOME/Games" 2>/dev/null | sort -hr >> "$TMP_DIRS"
    fi
    
    echo "    Verificando pastas visíveis no Home..."
    find "$REAL_HOME" -maxdepth 1 -not -path '*/.*' -not -path '*databases/datamysql*' -type d 2>/dev/null | xargs -I {} du -sh {} 2>/dev/null | sort -hr >> "$TMP_DIRS"
    
    find "$REAL_HOME/Downloads" "$REAL_HOME/Documents" -maxdepth 1 -type d 2>/dev/null | xargs -I {} du -sh {} 2>/dev/null | sort -hr | head -n 20 >> "$TMP_DIRS"
    
    sort -hr "$TMP_DIRS" | uniq > "${TMP_DIRS}.uniq"
    
    TOTAL_DIRS=$(wc -l < "${TMP_DIRS}.uniq")
    
    if [ "$TOTAL_DIRS" -eq 0 ]; then
        echo "    Nenhuma pasta grande encontrada."
    else
        echo "    Encontradas $TOTAL_DIRS pastas grandes/jogos."
        echo "    AVISO: Apagar pastas apaga TUDO dentro delas."
        
        while read -r line; do
            size=$(echo "$line" | awk '{print $1}')
            path=$(echo "$line" | cut -f2-)
            
            if [ -z "$path" ]; then continue; fi
            if [ "$path" == "." ] || [ "$path" == "$REAL_HOME" ]; then continue; fi
            if [ ! -d "$path" ]; then continue; fi

            echo "    Pasta: $path"
            echo "    Tamanho: $size"
            
            read -p "    >> EXCLUIR pasta completamenta? [s/n/q]: " -n 1 -r < /dev/tty
            echo ""
            
            case $REPLY in
                [sS])
                    rm -rf "$path"
                    echo "       ✓ Apagado."
                    ;;
                [qQ])
                    echo "       Parando lista de pastas..."
                    break
                    ;;
                *)
                    echo "       ⊘ Mantido."
                    ;;
            esac
            echo "---------------------------------------------------"
        done < "${TMP_DIRS}.uniq"
    fi
    rm -f "$TMP_DIRS" "${TMP_DIRS}.uniq"
else
    echo "    Busca de pastas pulada."
fi
