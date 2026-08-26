package model

type Family struct {
	ID       string
	Name     string
	Children int
	Notes    string
}

type Itinerary struct {
	ID          string
	FamilyID    string
	Title       string
	Destination string
	StartDate   string
	EndDate     string
	Status      string
}

type Activity struct {
	ID        string
	Itinerary string
	Day       int
	Title     string
	Category  string
	Duration  int
	ChildSafe bool
}

type ChatMessage struct {
	ID        string
	Itinerary string
	Sequence  int
	Sender    string
	Body      string
	Status    string
	Failure   string
}

type DateSuggestion struct {
	ID        string
	Itinerary string
	Suggested string
	Reason    string
	Accepted  bool
}

type Intake struct {
	Family      Family
	Itinerary   Itinerary
	Activities  []Activity
	Preferences []string
}

type ChatTranscript struct {
	Itinerary string
	Messages  []ChatMessage
}

func NewFamily(id, name string, children int, notes string) Family {
	return Family{ID: id, Name: name, Children: children, Notes: notes}
}

func NewItinerary(id, familyID, title, destination, start, end string) Itinerary {
	return Itinerary{ID: id, FamilyID: familyID, Title: title, Destination: destination, StartDate: start, EndDate: end, Status: "draft"}
}

func NewActivity(id, itinerary string, day int, title, category string, duration int, childSafe bool) Activity {
	return Activity{ID: id, Itinerary: itinerary, Day: day, Title: title, Category: category, Duration: duration, ChildSafe: childSafe}
}

func NewMessage(id, itinerary string, sequence int, sender, body string) ChatMessage {
	return ChatMessage{ID: id, Itinerary: itinerary, Sequence: sequence, Sender: sender, Body: body, Status: "pending"}
}

func NewSuggestion(id, itinerary, date, reason string) DateSuggestion {
	return DateSuggestion{ID: id, Itinerary: itinerary, Suggested: date, Reason: reason}
}
