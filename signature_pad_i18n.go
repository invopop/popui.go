package popui

type signaturePadLabels struct {
	Title       string
	DrawFirst   string
	ExportError string
	TabTyped    string
	TabDraw     string
	UseLabel    string
	ClearLabel  string
	CloseLabel  string
}

func spLabels(locale string) signaturePadLabels {
	switch locale {
	case "fr":
		return signaturePadLabels{
			Title:       "Signer",
			DrawFirst:   "Veuillez d'abord dessiner votre signature.",
			ExportError: "Impossible d'exporter la signature dans ce navigateur.",
			TabTyped:    "Cliquez pour signer",
			TabDraw:     "Dessinez votre signature",
			UseLabel:    "Utiliser la signature",
			ClearLabel:  "Effacer",
			CloseLabel:  "Fermer",
		}
	case "es":
		return signaturePadLabels{
			Title:       "Firmar",
			DrawFirst:   "Por favor, dibuje su firma primero.",
			ExportError: "No se puede exportar la firma en este navegador.",
			TabTyped:    "Haga clic para firmar",
			TabDraw:     "Dibuje su firma",
			UseLabel:    "Usar firma",
			ClearLabel:  "Borrar",
			CloseLabel:  "Cerrar",
		}
	default:
		return signaturePadLabels{
			Title:       "Sign",
			DrawFirst:   "Please draw your signature first.",
			ExportError: "Unable to export the signature in this browser.",
			TabTyped:    "Click to sign",
			TabDraw:     "Draw your signature",
			UseLabel:    "Use signature",
			ClearLabel:  "Clear",
			CloseLabel:  "Close",
		}
	}
}
