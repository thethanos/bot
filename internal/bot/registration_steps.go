package bot

import ma "multimessenger_bot/internal/messenger_adapter"

type Question struct {
	Text   string
	Answer string
	Field  string
}

var MasterQuestions = []*Question{
	{Text: "Name?"},
	{Text: "City?"},
	{Text: "Service?"},
}

type Master struct {
	StepBase
}

func (m *Master) Request(msg *ma.Message) *ma.Message {
	m.inProgress = true
	return &ma.Message{Text: "Register?", UserData: msg.UserData, Type: msg.Type}
}

func (m *Master) ProcessResponse(msg *ma.Message) (*ma.Message, int) {
	m.inProgress = false
	return nil, RegistrationStep
}

func (m *Master) IsInProgress() bool {
	return m.inProgress
}

type Registration struct {
	StepBase
	Questions []*Question
}

func (r *Registration) Request(msg *ma.Message) *ma.Message {
	r.inProgress = true
	return &ma.Message{Text: r.Questions[r.State.cursor].Text, UserData: msg.UserData, Type: msg.Type}
}

func (r *Registration) ProcessResponse(msg *ma.Message) (*ma.Message, int) {
	r.inProgress = false
	r.Questions[r.State.cursor].Answer = msg.Text
	r.State.cursor = r.State.cursor + 1
	if r.State.cursor >= len(r.Questions) {
		r.State.cursor = 0
		return nil, RegistrationFinalStep
	}
	return nil, RegistrationStep
}

func (r *Registration) IsInProgress() bool {
	return r.inProgress
}

type RegistrationFinal struct {
	StepBase
}

func (r *RegistrationFinal) Request(msg *ma.Message) *ma.Message {
	r.inProgress = true
	return nil
}

func (r *RegistrationFinal) ProcessResponse(msg *ma.Message) (*ma.Message, int) {
	return nil, MainMenuStep
}

func (r *RegistrationFinal) IsInProgress() bool {
	return r.inProgress
}
