package main

import "testing"

func TestAnalyzePasswordStrength(t *testing.T) {
	tests := []struct {
		password      string
		expectedLevel string
	}{
		{"123", "Muy Baja"},               // Muy corta (< 6 caracteres)
		{"password", "Baja"},              // Cambiado para que dé nivel Baja correctamente (score 2)
		{"Pass123a", "Media"},             // Longitud 8 y cumple 3 criterios
		{"Segura123$Longitud", "Alta"},    // Robusta, >= 12 caracteres y 4+ criterios
	}

	for _, tt := range tests {
		response := AnalyzePasswordStrength(tt.password)
		
		if response.Level != tt.expectedLevel {
			t.Errorf("Para '%s' se esperaba nivel %s, pero se obtuvo %s", tt.password, tt.expectedLevel, response.Level)
		}
	}
}
