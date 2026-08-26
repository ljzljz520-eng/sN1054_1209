package model

func ValidFamily(f Family) bool {
	if f.ID == "" || f.Name == "" {
		return false
	}
	if f.Children < 0 || f.Children > 12 {
		return false
	}
	return true
}

func ValidItinerary(i Itinerary) bool {
	if i.ID == "" || i.FamilyID == "" || i.Destination == "" {
		return false
	}
	if i.StartDate == "" || i.EndDate == "" {
		return false
	}
	return i.StartDate <= i.EndDate
}

func ValidActivity(a Activity) bool {
	if a.ID == "" || a.Itinerary == "" || a.Title == "" {
		return false
	}
	if a.Day < 1 || a.Duration <= 0 {
		return false
	}
	return true
}

func ValidMessage(m ChatMessage) bool {
	if m.ID == "" || m.Itinerary == "" || m.Body == "" {
		return false
	}
	if m.Sequence < 1 || m.Sender == "" {
		return false
	}
	return true
}

func NormalizeStatus(status string) string {
	if status == "" {
		return "pending"
	}
	switch status {
	case "pending", "sent", "failed", "accepted", "rejected":
		return status
	default:
		return "pending"
	}
}

func IsFailure(status string) bool {
	return status == "failed" || status == "rejected"
}

func PreferenceAllowed(value string) bool {
	if value == "" {
		return false
	}
	switch value {
	case "nap-friendly", "short-transfers", "indoor-backup", "early-dinner", "stroller-access":
		return true
	default:
		return false
	}
}
