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

type Question struct {
	Text  string
	Field string
}

type RegistrationFinal struct {
	StepBase
}

func (r *RegistrationFinal) Request(msg *ma.Message) *ma.Message {
	r.logger.Infof("RegistrationFinal step is sending request")
	r.inProgress = true
	data := FormatMapToString(r.state.RawInput)
	if msg.Source == ma.TELEGRAM {
		rows := make([][]tgbotapi.KeyboardButton, 0)
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Да"}})
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Нет"}})
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}
		return ma.NewTextMessage(fmt.Sprintf("%s\nПодтвердить регистрацию?", data), msg, keyboard, false)
	}
	return ma.NewTextMessage(fmt.Sprintf("%s\nПодтвердить регистрацию?\n1. Да\n2. Нет", data), msg, nil, true)
}

func (r *RegistrationFinal) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	r.logger.Infof("RegistrationFinal step is processing response")
	r.inProgress = false
	userAnswer := strings.ToLower(msg.Text)
	if userAnswer == "да" || userAnswer == "1" {
		r.dbAdapter.SaveMaster(r.state)
		r.state.Reset()
		return ma.NewTextMessage("Регистрация прошла успешно!", msg, nil, true), MainMenuRequestStep
	}
	r.state.Reset()
	return nil, MainMenuRequestStep
}

func FormatMapToString(data map[string]string) string {
	res := ""
	for key, val := range data {
		res += fmt.Sprintf("%s: %s\n", key, val)
	}
	return res
}

type Downloader interface {
	DownloadFile(ma.FileType, *ma.Message) []byte
}

type ImageUpload struct {
	StepBase
	downloader Downloader
}

func (i *ImageUpload) Request(msg *ma.Message) *ma.Message {
	i.logger.Info("ImageUpload step is sending request")
	i.inProgress = true
	text := "Добавить фото"
	if msg.Source == ma.TELEGRAM {
		rows := make([][]tgbotapi.KeyboardButton, 0)
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Далее"}})
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Назад"}})
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Главное меню"}})
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}
		return ma.NewTextMessage(text, msg, keyboard, false)
	}
	return ma.NewTextMessage(text, msg, nil, false)
}

func (i *ImageUpload) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	i.logger.Info("ImageUpload step is processing response")
	userAnswer := strings.ToLower(msg.Text)
	if userAnswer == "далее" {
		if err := i.dbAdapter.SaveMasterPreview(i.state.Master); err != nil {
			return ma.NewTextMessage("Не удалось сохранить данные для предпросмотра", msg, nil, false), EmptyStep
		}
		return nil, AddMasterFinalStep
	}
	if userAnswer == "назад" {
		return nil, PreviousStep
	}
	if userAnswer == "главное меню" {
		return nil, MainMenuStep
	}
	file := i.downloader.DownloadFile(ma.PHOTO, msg)
	if len(file) == 0 {
		return ma.NewTextMessage("Не удалось загрузить изображение", msg, nil, false), EmptyStep
	}

	image := parsers.SaveFile(i.state.Master.ID, "./images", "jpeg", file)
	i.state.Master.Images = append(i.state.Master.Images, fmt.Sprintf("https://bot-dev-domain.com/pages/images/%s/%s", i.state.Master.ID, image))
	return nil, EmptyStep
}

type AddMasterFinal struct {
	StepBase
}

func (a *AddMasterFinal) Request(msg *ma.Message) *ma.Message {
	a.logger.Info("AddMasterFinal step is sending request")
	a.inProgress = true
	text := "Завершающий этап"
	if msg.Source == ma.TELEGRAM {
		rows := make([][]tgbotapi.KeyboardButton, 0)
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Предпросмотр", WebApp: &tgbotapi.WebAppInfo{Url: fmt.Sprintf("https://bot-dev-domain.com/masters/preview?master=%s", a.state.Master.ID)}}})
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Сохранить анкету"}})
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Назад"}})
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Главное меню"}})
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}
		return ma.NewTextMessage(text, msg, keyboard, false)
	}
	return ma.NewTextMessage(text, msg, nil, false)
}

func (a *AddMasterFinal) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	a.logger.Info("AddMasterFinal step is processing response")
	userAnswer := strings.ToLower(msg.Text)
	if userAnswer == "назад" {
		return nil, PreviousStep
	}
	if userAnswer == "главное меню" {
		return nil, MainMenuStep
	}
	return ma.NewTextMessage("Анкета сохранена", msg, nil, false), AdminStep
}


type AddMaster struct {
	StepBase
	downloader Downloader
}

func (a *AddMaster) Request(msg *ma.Message) *ma.Message {
	a.logger.Info("AddMaster step is sending request")
	a.inProgress = true
	text := "Загрузите данные мастера в формате xlsx"
	if msg.Source == ma.TELEGRAM {
		rows := make([][]tgbotapi.KeyboardButton, 0)
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Назад"}})
		rows = append(rows, []tgbotapi.KeyboardButton{{Text: "Главное меню"}})
		keyboard := &tgbotapi.ReplyKeyboardMarkup{Keyboard: rows, ResizeKeyboard: true}
		return ma.NewTextMessage(text, msg, keyboard, false)
	}
	return ma.NewTextMessage(text, msg, nil, true)
}

func (a *AddMaster) ProcessResponse(msg *ma.Message) (*ma.Message, StepType) {
	a.logger.Info("AddMaster step is processing response")
	userAnswer := strings.ToLower(msg.Text)
	if userAnswer == "назад" {
		return nil, PreviousStep
	}
	if userAnswer == "главное меню" {
		return nil, MainMenuStep
	}

	file := a.downloader.DownloadFile(ma.DOCUMENT, msg)

	master, err := parsers.ParseMasterData(file)
	if err != nil {
		return ma.NewTextMessage("Не удалось распарсить данные мастера, проверьте правильность ввода и попробуйте еще раз", msg, nil, false), EmptyStep
	}
	a.state.Master = master
	a.inProgress = false
	return nil, ImageUploadStep
}

*/
