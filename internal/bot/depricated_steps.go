package bot

/*

type MasterSelection struct {
	StepBase
	masters []*entities.Master
}

func (m *MasterSelection) Request(msg *ma.Message) *ma.Message {
	m.logger.Infof("MasterSelection step is sending request")
	m.inProgress = true
	masters, _ := m.DbAdapter.GetMasters(m.State.City.ID, m.State.Service.ID)

	if msg.Source == ma.TELEGRAM {

		rows := make([][]tgbotapi.KeyboardButton, len(masters)+1)
		for idx, master := range masters {
			rows[idx] = make([]tgbotapi.KeyboardButton, 0)
			rows[idx] = append(rows[idx], tgbotapi.KeyboardButton{Text: master.Name})
		}
		rows[len(masters)] = make([]tgbotapi.KeyboardButton, 0)
		rows[len(masters)] = append(rows[len(masters)], tgbotapi.KeyboardButton{Text: "Назад"})
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}

		if len(masters) == 0 {
			return ma.NewMessage("По вашему запросу ничего не найдено", ma.REGULAR, msg, keyboard, nil)
		}

		m.masters = masters
		return ma.NewMessage(" Выберите мастера", ma.REGULAR, msg, keyboard, nil)
	}

	text := ""
	for idx, master := range masters {
		text += fmt.Sprintf("%d. %s", idx+1, master.Name)
	}

	m.masters = masters
	return ma.NewMessage(text, ma.REGULAR, msg, nil, nil)
}

func (m *MasterSelection) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	if msg.Type == ma.CALLBACK {
		return nil, EmptyStep
	}
	m.logger.Infof("MasterSelection step is processing response")
	m.inProgress = false

	userAnswer := strings.ToLower(msg.Text)
	if userAnswer == "назад" || userAnswer == fmt.Sprintf("%d", len(m.masters)+1) {
		m.logger.Infof("Next step is PreviousStep")
		return nil, PreviousStep
	}
	for idx, master := range m.masters {
		if userAnswer == strings.ToLower(master.Name) || userAnswer == fmt.Sprintf("%d", idx+1) {
			m.State.Master = master
			m.logger.Infof("Next step is FinalStep")
			return nil, FinalStep
		}
	}

	m.inProgress = true
	m.logger.Infof("Next step is EmptyStep")
	return ma.NewMessage("Пожалуйста выберите ответ из списка.", ma.REGULAR, msg, nil, nil), EmptyStep
}

type Final struct {
	StepBase
}

func (f *Final) Request(msg *ma.Message) *ma.Message {
	f.logger.Infof("Final step is sending request")
	f.inProgress = true
	text := fmt.Sprintf("Ваша запись\nУслуга: %s\nГород: %s\nМастер: %s\n\nПодтвердить?",
		f.State.Service.Name,
		f.State.City.Name,
		f.State.Master.Name,
	)

	if msg.Source == ma.TELEGRAM {
		rows := make([][]tgbotapi.KeyboardButton, 2)
		rows[0] = []tgbotapi.KeyboardButton{{Text: "Да"}}
		rows[1] = []tgbotapi.KeyboardButton{{Text: "Нет"}}
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}
		return ma.NewMessage(text, ma.REGULAR, msg, keyboard, nil)
	}
	return ma.NewMessage(fmt.Sprintf("%s\n1. Да\n2. Нет", text), ma.REGULAR, msg, nil, nil)
}

func (f *Final) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	if msg.Type == ma.CALLBACK {
		return nil, EmptyStep
	}
	f.logger.Infof("Final step is processing response")
	f.inProgress = false

	switch msg.Text {
	case "Да":
		f.State.Reset()
		f.logger.Infof("Next step is MainMenuRequestStep")
		return ma.NewMessage("Запись завершена", ma.REGULAR, msg, nil, nil), MainMenuRequestStep
	case "Нет":
		f.State.Reset()
		f.logger.Infof("Next step is MainMenuRequestStep")
		return ma.NewMessage("Запись отменена", ma.REGULAR, msg, nil, nil), MainMenuRequestStep
	default:
		f.inProgress = true
		f.logger.Infof("Next step is EmptyStep")
		return ma.NewMessage("Пожалуйста выберите ответ из списка.", ma.REGULAR, msg, nil, nil), EmptyStep
	}
}

type CityPromptStepMode interface {
	Text() string
	Buttons() *tgbotapi.ReplyKeyboardMarkup
	NextStep() StepType
}

type BaseCityPromptMode struct {
}

func (b *BaseCityPromptMode) Text() string {
	return "Введите город"
}

func (b *BaseCityPromptMode) Buttons() *tgbotapi.ReplyKeyboardMarkup {
	rows := make([][]tgbotapi.KeyboardButton, 0)
	rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Назад"}})
	rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Главное меню"}})
	return &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}
}

func (b *BaseCityPromptMode) NextStep() StepType {
	return EmptyStep
}

type MainMenuCityPromptMode struct {
	BaseCityPromptMode
}

func (m *MainMenuCityPromptMode) Buttons() *tgbotapi.ReplyKeyboardMarkup {
	rows := make([][]tgbotapi.KeyboardButton, 0)
	rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Главное меню"}})
	return &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}
}

func (m *MainMenuCityPromptMode) NextStep() StepType {
	return ServiceCategorySelectionStep
}

type RegistrationCityPromptMode struct {
	BaseCityPromptMode
}

func (r *RegistrationCityPromptMode) NextStep() StepType {
	return MasterServiceCategorySecletionStep
}

type CityPrompt struct {
	StepBase
	mode CityPromptStepMode
}

func (c *CityPrompt) Request(msg *ma.Message) *ma.Message {
	c.logger.Infof("CityPrompt step is sending request")
	c.inProgress = true

	if msg.Source == ma.TELEGRAM {
		return ma.NewTextMessage(c.mode.Text(), msg, c.mode.Buttons(), false)
	}

	return ma.NewTextMessage(fmt.Sprintf("%s\n1. Назад\n2. Главное меню", c.mode.Text()), msg, nil, true)
}

func (c *CityPrompt) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	c.logger.Infof("CityPrompt step is processing response")
	c.inProgress = false

	userAnswer := strings.ToLower(msg.Text)
	if userAnswer == "назад" {
		return nil, PreviousStep
	}
	if userAnswer == "главное меню" {
		return nil, MainMenuStep
	}

	city, err := c.dbAdapter.GetCity(msg.Text)
	if err != nil {
		c.inProgress = true
		c.logger.Infof("Next step is CityPromptStep")
		return ma.NewTextMessage(fmt.Sprintf("По запросу %s ничего не найдено", msg.Text), msg, nil, false), CityPromptStep
	}
	c.state.City = city
	c.logger.Infof("Next step is %s", getStepTypeName(c.mode.NextStep()))
	return nil, c.mode.NextStep()
}

func (c *CityPrompt) Reset() {
	c.state.City = nil
}
*/
