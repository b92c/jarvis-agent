package main

import (
	"testing"
	"time"
)

func TestGetNextRunTime(t *testing.T) {
	location, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		t.Fatalf("Erro ao carregar timezone: %v", err)
	}

	tests := []struct {
		name         string
		now          time.Time
		expectedTime time.Time
		expectedType runType
	}{
		{
			name:         "antes das 12h - retorna 12h mesmo dia",
			now:          time.Date(2024, 1, 15, 10, 0, 0, 0, location),
			expectedTime: time.Date(2024, 1, 15, 12, 0, 0, 0, location),
			expectedType: runDaily,
		},
		{
			name:         "entre 12h e 18h - retorna 18h mesmo dia",
			now:          time.Date(2024, 1, 15, 14, 30, 0, 0, location),
			expectedTime: time.Date(2024, 1, 15, 18, 0, 0, 0, location),
			expectedType: runDaily,
		},
		{
			name:         "depois das 18h - retorna 12h do dia seguinte",
			now:          time.Date(2024, 1, 15, 19, 0, 0, 0, location),
			expectedTime: time.Date(2024, 1, 16, 12, 0, 0, 0, location),
			expectedType: runDaily,
		},
		{
			name:         "exatamente às 12h - retorna 18h mesmo dia",
			now:          time.Date(2024, 1, 15, 12, 0, 0, 0, location),
			expectedTime: time.Date(2024, 1, 15, 18, 0, 0, 0, location),
			expectedType: runDaily,
		},
		{
			name:         "exatamente às 18h - retorna 12h do dia seguinte",
			now:          time.Date(2024, 1, 15, 18, 0, 0, 0, location),
			expectedTime: time.Date(2024, 1, 16, 12, 0, 0, 0, location),
			expectedType: runDaily,
		},
		{
			name:         "último dia do mês às 18h - retorna relatório mensal",
			now:          time.Date(2024, 1, 31, 17, 0, 0, 0, location),
			expectedTime: time.Date(2024, 1, 31, 18, 0, 0, 0, location),
			expectedType: runMonthly,
		},
		{
			name:         "último dia do mês antes das 12h - retorna 12h do mesmo dia (daily)",
			now:          time.Date(2024, 1, 31, 10, 0, 0, 0, location),
			expectedTime: time.Date(2024, 1, 31, 12, 0, 0, 0, location),
			expectedType: runDaily,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getNextRunTime(tt.now)
			if !result.time.Equal(tt.expectedTime) {
				t.Errorf("getNextRunTime(%v).time = %v, want %v", tt.now, result.time, tt.expectedTime)
			}
			if result.runType != tt.expectedType {
				t.Errorf("getNextRunTime(%v).runType = %v, want %v", tt.now, result.runType, tt.expectedType)
			}
			if result.time.Location().String() != location.String() {
				t.Errorf("Timezone retornado = %v, want %v", result.time.Location(), location)
			}
		})
	}
}

func TestSaoPauloLocation(t *testing.T) {
	if saoPauloLocation == nil {
		t.Fatal("saoPauloLocation não foi inicializado")
	}
	if saoPauloLocation.String() != "America/Sao_Paulo" {
		t.Errorf("saoPauloLocation = %v, want America/Sao_Paulo", saoPauloLocation)
	}
}
