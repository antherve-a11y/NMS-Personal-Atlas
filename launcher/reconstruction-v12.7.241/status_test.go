package main

import "testing"

func TestPhaseParsing(t *testing.T) {
	s := phaseState{}
	for _, line := range []string{
		"@ATLAS_UI\tDATA\tInitialisation", "@ATLAS_UI\tCHECK\tContrôle", "@ATLAS_UI\tINSTALL\tTéléchargement", "@ATLAS_UI\tVERIFY\tVérification", "@ATLAS_UI\tREADY\tPrêt",
	} {
		s = parseStatusLine(line, s)
	}
	if s.Phase != 4 || !s.Done || s.Text != "Prêt" {
		t.Fatalf("%+v", s)
	}
}
func TestFallback(t *testing.T) {
	s := parseStatusLine("[ATLAS] Installation de Streamlit, pandas, NumPy...", phaseState{})
	if s.Phase != 3 {
		t.Fatalf("%+v", s)
	}
}
