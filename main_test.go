package main

import "testing"

func TestCheckPasswordStrength(t *testing.T) {
	tests := []struct {
		password      string
		expectedLevel string
	}{
		{"123", "Baja"},                  // Muy corta
		{"password123", "Baja"},           // Larga pero sin mayúsculas/especiales
		{"Pass123a", "Media"},            // MODIFICADO: longitud 8 y cumple 3 criterios
		{"Segura123$Longitud", "Alta"},    // Robusta, >12 caracteres y todos los criterios
	}

	for _, tt := range tests {
		level, _ := CheckPasswordStrength(tt.password)
		if level != tt.expectedLevel {
			t.Errorf("Para '%s' se esperaba nivel %s, pero se obtuvo %s", tt.password, tt.expectedLevel, level)
		}
	}
}
