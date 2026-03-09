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
		name     string
		now      time.Time
		expected time.Time
	}{
		{
			name:     "antes das 12h - retorna 12h mesmo dia",
			now:      time.Date(2024, 1, 15, 10, 0, 0, 0, location),
			expected: time.Date(2024, 1, 15, 12, 0, 0, 0, location),
		},
		{
			name:     "entre 12h e 18h - retorna 18h mesmo dia",
			now:      time.Date(2024, 1, 15, 14, 30, 0, 0, location),
			expected: time.Date(2024, 1, 15, 18, 0, 0, 0, location),
		},
		{
			name:     "depois das 18h - retorna 12h do dia seguinte",
			now:      time.Date(2024, 1, 15, 19, 0, 0, 0, location),
			expected: time.Date(2024, 1, 16, 12, 0, 0, 0, location),
		},
		{
			name:     "exatamente às 12h - retorna 18h mesmo dia",
			now:      time.Date(2024, 1, 15, 12, 0, 0, 0, location),
			expected: time.Date(2024, 1, 15, 18, 0, 0, 0, location),
		},
		{
			name:     "exatamente às 18h - retorna 12h do dia seguinte",
			now:      time.Date(2024, 1, 15, 18, 0, 0, 0, location),
			expected: time.Date(2024, 1, 16, 12, 0, 0, 0, location),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := getNextRunTime(tt.now)
			if !result.Equal(tt.expected) {
				t.Errorf("getNextRunTime(%v) = %v, want %v", tt.now, result, tt.expected)
			}
			if result.Location().String() != location.String() {
				t.Errorf("Timezone retornado = %v, want %v", result.Location(), location)
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
