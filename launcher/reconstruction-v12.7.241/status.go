package main

import "strings"

type phaseState struct {
	Phase int
	Text  string
	Done  bool
	Error string
}

// Interpret procedural output. Prefer tagged events; support older Atlas logs.
func parseStatusLine(s string, state phaseState) phaseState {
	s = strings.TrimSpace(s)
	if strings.HasPrefix(s, "@ATLAS_UI\t") {
		fields := strings.SplitN(s, "\t", 3)
		if len(fields) == 3 {
			switch fields[1] {
			case "DATA":
				state.Phase = 1
				state.Text = fields[2]
			case "CHECK":
				state.Phase = 2
				state.Text = fields[2]
			case "INSTALL":
				state.Phase = 3
				state.Text = fields[2]
			case "VERIFY":
				state.Phase = 4
				state.Text = fields[2]
			case "START":
				state.Phase = 4
				state.Text = fields[2]
			case "READY":
				state.Phase = 4
				state.Text = fields[2]
				state.Done = true
			case "ERROR":
				state.Error = fields[2]
			}
		}
		return state
	}
	lower := strings.ToLower(s)
	switch {
	case strings.Contains(lower, "donnees locales"), strings.Contains(lower, "données locales"):
		if state.Phase < 1 {
			state.Phase = 1
			state.Text = "Préparation des données personnelles"
		}
	case strings.Contains(lower, "premiere utilisation"):
		state.Phase = 2
		state.Text = "Vérification de l'environnement Atlas"
	case strings.Contains(lower, "installation de streamlit"):
		state.Phase = 3
		state.Text = "Installation des bibliothèques nécessaires"
	case strings.Contains(lower, "bibliotheques installees"):
		state.Phase = 4
		state.Text = "Vérification terminée. Ouverture d'Atlas…"
	case strings.Contains(lower, "demarrage local"):
		state.Phase = 4
		state.Text = "Démarrage du serveur Atlas…"
	case strings.Contains(lower, "installation automatique incomplete"), strings.Contains(lower, "une erreur"):
		state.Error = s
	}
	return state
}
