#!/bin/bash

# Script para testar carregamento de bateria no Linux
# Autor: Script de teste de bateria

echo "=========================================="
echo "  TESTE DE CARREGAMENTO DE BATERIA"
echo "=========================================="
echo ""

# Cores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Função para ler valor de arquivo sysfs
read_sysfs() {
    local file=$1
    if [ -f "$file" ]; then
        cat "$file" 2>/dev/null
    else
        echo "N/A"
    fi
}

# Verificar se os diretórios existem
BAT_DIR="/sys/class/power_supply/BAT0"
AC_DIR="/sys/class/power_supply/ADP0"

if [ ! -d "$BAT_DIR" ]; then
    echo -e "${RED}Erro: Diretório da bateria não encontrado!${NC}"
    exit 1
fi

echo -e "${BLUE}=== INFORMAÇÕES DA BATERIA ===${NC}"
echo ""

# Status da bateria
STATUS=$(read_sysfs "$BAT_DIR/status")
CAPACITY=$(read_sysfs "$BAT_DIR/capacity")
CAPACITY_LEVEL=$(read_sysfs "$BAT_DIR/capacity_level")
ENERGY_NOW=$(read_sysfs "$BAT_DIR/energy_now")
ENERGY_FULL=$(read_sysfs "$BAT_DIR/energy_full")
ENERGY_FULL_DESIGN=$(read_sysfs "$BAT_DIR/energy_full_design")
POWER_NOW=$(read_sysfs "$BAT_DIR/power_now")
VOLTAGE_NOW=$(read_sysfs "$BAT_DIR/voltage_now")
CYCLE_COUNT=$(read_sysfs "$BAT_DIR/cycle_count")
TECHNOLOGY=$(read_sysfs "$BAT_DIR/technology")
MANUFACTURER=$(read_sysfs "$BAT_DIR/manufacturer")
MODEL_NAME=$(read_sysfs "$BAT_DIR/model_name")

# Status do carregador
if [ -d "$AC_DIR" ]; then
    AC_ONLINE=$(read_sysfs "$AC_DIR/online")
else
    AC_ONLINE="N/A"
fi

# Exibir informações
echo -e "Status: ${YELLOW}$STATUS${NC}"
echo -e "Capacidade: ${GREEN}$CAPACITY%${NC}"
if [ "$CAPACITY_LEVEL" != "N/A" ] && [ "$CAPACITY_LEVEL" != "Unknown" ]; then
    echo "Nível de capacidade: $CAPACITY_LEVEL"
fi

echo ""
echo -e "${BLUE}=== STATUS DO CARREGADOR ===${NC}"
if [ "$AC_ONLINE" = "1" ]; then
    echo -e "Carregador: ${GREEN}Conectado${NC}"
elif [ "$AC_ONLINE" = "0" ]; then
    echo -e "Carregador: ${RED}Desconectado${NC}"
else
    echo "Carregador: $AC_ONLINE"
fi

echo ""
echo -e "${BLUE}=== DETALHES TÉCNICOS ===${NC}"

# Converter valores (geralmente em micro-unidades)
if [ "$ENERGY_NOW" != "N/A" ] && [ "$ENERGY_NOW" != "" ]; then
    ENERGY_NOW_WH=$(echo "scale=2; $ENERGY_NOW / 1000000" | bc 2>/dev/null || echo "$ENERGY_NOW")
    echo "Energia atual: ${ENERGY_NOW_WH} Wh"
fi

if [ "$ENERGY_FULL" != "N/A" ] && [ "$ENERGY_FULL" != "" ]; then
    ENERGY_FULL_WH=$(echo "scale=2; $ENERGY_FULL / 1000000" | bc 2>/dev/null || echo "$ENERGY_FULL")
    echo "Energia total: ${ENERGY_FULL_WH} Wh"
fi

if [ "$ENERGY_FULL_DESIGN" != "N/A" ] && [ "$ENERGY_FULL_DESIGN" != "" ]; then
    ENERGY_DESIGN_WH=$(echo "scale=2; $ENERGY_FULL_DESIGN / 1000000" | bc 2>/dev/null || echo "$ENERGY_FULL_DESIGN")
    echo "Energia de design: ${ENERGY_DESIGN_WH} Wh"
    
    # Calcular saúde da bateria
    if [ "$ENERGY_FULL" != "N/A" ] && [ "$ENERGY_FULL" != "" ] && [ "$ENERGY_FULL_DESIGN" != "N/A" ] && [ "$ENERGY_FULL_DESIGN" != "" ]; then
        HEALTH=$(echo "scale=1; ($ENERGY_FULL * 100) / $ENERGY_FULL_DESIGN" | bc 2>/dev/null)
        if [ ! -z "$HEALTH" ]; then
            echo -e "Saúde da bateria: ${GREEN}${HEALTH}%${NC}"
        fi
    fi
fi

if [ "$POWER_NOW" != "N/A" ] && [ "$POWER_NOW" != "" ]; then
    POWER_NOW_W=$(echo "scale=2; $POWER_NOW / 1000000" | bc 2>/dev/null || echo "$POWER_NOW")
    if [ "$STATUS" = "Charging" ]; then
        echo -e "Potência de carregamento: ${GREEN}+${POWER_NOW_W} W${NC}"
    elif [ "$STATUS" = "Discharging" ]; then
        echo -e "Potência de descarga: ${RED}-${POWER_NOW_W} W${NC}"
    else
        echo "Potência: ${POWER_NOW_W} W"
    fi
fi

if [ "$VOLTAGE_NOW" != "N/A" ] && [ "$VOLTAGE_NOW" != "" ]; then
    VOLTAGE_NOW_V=$(echo "scale=2; $VOLTAGE_NOW / 1000000" | bc 2>/dev/null || echo "$VOLTAGE_NOW")
    echo "Tensão: ${VOLTAGE_NOW_V} V"
fi

if [ "$CYCLE_COUNT" != "N/A" ] && [ "$CYCLE_COUNT" != "" ]; then
    echo "Ciclos de carga: $CYCLE_COUNT"
fi

if [ "$TECHNOLOGY" != "N/A" ] && [ "$TECHNOLOGY" != "" ]; then
    echo "Tecnologia: $TECHNOLOGY"
fi

if [ "$MANUFACTURER" != "N/A" ] && [ "$MANUFACTURER" != "" ]; then
    echo "Fabricante: $MANUFACTURER"
fi

if [ "$MODEL_NAME" != "N/A" ] && [ "$MODEL_NAME" != "" ]; then
    echo "Modelo: $MODEL_NAME"
fi

echo ""
echo -e "${BLUE}=== TESTE DE CARREGAMENTO ===${NC}"

# Verificar se está carregando
if [ "$STATUS" = "Charging" ] && [ "$AC_ONLINE" = "1" ]; then
    echo -e "${GREEN}✓ Bateria está carregando normalmente${NC}"
    
    # Estimar tempo de carregamento
    if [ "$POWER_NOW" != "N/A" ] && [ "$POWER_NOW" != "" ] && [ "$ENERGY_NOW" != "N/A" ] && [ "$ENERGY_NOW" != "" ] && [ "$ENERGY_FULL" != "N/A" ] && [ "$ENERGY_FULL" != "" ]; then
        ENERGY_REMAINING=$(echo "$ENERGY_FULL - $ENERGY_NOW" | bc 2>/dev/null)
        if [ ! -z "$ENERGY_REMAINING" ] && [ "$ENERGY_REMAINING" != "0" ] && [ "$POWER_NOW" != "0" ]; then
            HOURS=$(echo "scale=1; $ENERGY_REMAINING / $POWER_NOW" | bc 2>/dev/null)
            if [ ! -z "$HOURS" ]; then
                HOURS_INT=$(echo "$HOURS" | cut -d. -f1)
                MINUTES=$(echo "scale=0; ($HOURS - $HOURS_INT) * 60" | bc 2>/dev/null)
                echo "Tempo estimado para carregar: ~${HOURS_INT}h ${MINUTES}m"
            fi
        fi
    fi
elif [ "$STATUS" = "Full" ] && [ "$AC_ONLINE" = "1" ]; then
    echo -e "${GREEN}✓ Bateria totalmente carregada${NC}"
elif [ "$STATUS" = "Discharging" ] && [ "$AC_ONLINE" = "0" ]; then
    echo -e "${YELLOW}⚠ Bateria descarregando (carregador desconectado)${NC}"
elif [ "$STATUS" = "Discharging" ] && [ "$AC_ONLINE" = "1" ]; then
    echo -e "${RED}⚠ ATENÇÃO: Carregador conectado mas bateria descarregando!${NC}"
    echo "   Isso pode indicar um problema com o carregador ou a bateria."
elif [ "$AC_ONLINE" = "1" ] && [ "$STATUS" != "Charging" ] && [ "$STATUS" != "Full" ]; then
    echo -e "${YELLOW}⚠ Carregador conectado mas status da bateria: $STATUS${NC}"
else
    echo -e "${YELLOW}Status: $STATUS${NC}"
fi

echo ""
echo "=========================================="

# Modo de monitoramento contínuo
if [ "$1" = "--monitor" ] || [ "$1" = "-m" ]; then
    echo ""
    echo -e "${BLUE}Modo de monitoramento ativado (Ctrl+C para sair)${NC}"
    echo ""
    
    while true; do
        clear
        echo "=========================================="
        echo "  MONITORAMENTO DE BATERIA"
        echo "  $(date '+%Y-%m-%d %H:%M:%S')"
        echo "=========================================="
        echo ""
        
        STATUS=$(read_sysfs "$BAT_DIR/status")
        CAPACITY=$(read_sysfs "$BAT_DIR/capacity")
        AC_ONLINE=$(read_sysfs "$AC_DIR/online")
        POWER_NOW=$(read_sysfs "$BAT_DIR/power_now")
        
        echo -e "Status: ${YELLOW}$STATUS${NC}"
        echo -e "Capacidade: ${GREEN}$CAPACITY%${NC}"
        
        if [ "$AC_ONLINE" = "1" ]; then
            echo -e "Carregador: ${GREEN}Conectado${NC}"
        else
            echo -e "Carregador: ${RED}Desconectado${NC}"
        fi
        
        if [ "$POWER_NOW" != "N/A" ] && [ "$POWER_NOW" != "" ]; then
            POWER_NOW_W=$(echo "scale=2; $POWER_NOW / 1000000" | bc 2>/dev/null || echo "$POWER_NOW")
            if [ "$STATUS" = "Charging" ]; then
                echo -e "Potência: ${GREEN}+${POWER_NOW_W} W${NC}"
            elif [ "$STATUS" = "Discharging" ]; then
                echo -e "Potência: ${RED}-${POWER_NOW_W} W${NC}"
            fi
        fi
        
        echo ""
        echo "Pressione Ctrl+C para sair..."
        sleep 2
    done
fi
