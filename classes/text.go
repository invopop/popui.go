package classes

// TextSize returns the classes of a text size from the props.TextSize
// constants, or an empty string for any other value. The classes are
// written out in full so Tailwind finds them here.
func TextSize(size string) string {
	switch size {
	case "sm":
		return "text-sm"
	case "base":
		return "text-base"
	case "lg":
		return "text-lg"
	case "xl":
		return "text-xl"
	}
	return ""
}

// TextTone returns the colour of a text tone from the props.TextTone
// constants, or an empty string for any other value.
func TextTone(tone string) string {
	switch tone {
	case "default":
		return "text-foreground"
	case "secondary":
		return "text-foreground-default-secondary"
	case "tertiary":
		return "text-foreground-default-tertiary"
	}
	return ""
}
