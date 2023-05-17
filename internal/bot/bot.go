package bot

import (
	"fmt"
	"multimessenger_bot/internal/db_adapter"
	"multimessenger_bot/internal/entities"
	ma "multimessenger_bot/internal/messenger_adapter"
	"time"

	"go.uber.org/zap"
)

type UserSession struct {
	CurrentStep         Step
	СurrentCallbackStep Step
	PrevSteps           StepStack
	PrevCallBackSteps   StepStack
	State               *entities.UserState
}

type Bot struct {
	logger       *zap.SugaredLogger
	clients      map[ma.MessageSource]ma.ClientInterface
	userSessions map[string]*UserSession
	recvMsgChan  chan *ma.Message
	sendMsgChan  chan *ma.Message
	dbAdapter    *db_adapter.DbAdapter
}

func NewBot(logger *zap.SugaredLogger, clientArray []ma.ClientInterface, dbAdpter *db_adapter.DbAdapter, recvMsgChan chan *ma.Message) (*Bot, error) {

	clients := make(map[ma.MessageSource]ma.ClientInterface)
	for _, client := range clientArray {
		clients[client.GetType()] = client
	}

	userSessions := make(map[string]*UserSession)
	sendMsgChan := make(chan *ma.Message)

	bot := &Bot{
		logger:       logger,
		clients:      clients,
		userSessions: userSessions,
		recvMsgChan:  recvMsgChan,
		dbAdapter:    dbAdpter,
		sendMsgChan:  sendMsgChan,
	}

	return bot, nil
}

func (b *Bot) Run() {

	for _, client := range b.clients {
		client.Connect()
	}

	go func() {
		for msg := range b.recvMsgChan {
			if _, exists := b.userSessions[msg.UserID]; !exists {
				state := &entities.UserState{Cursor: 0, RawInput: make(map[string]string)}
				b.userSessions[msg.UserID] = &UserSession{
					State:       state,
					CurrentStep: b.createStep(MainMenuStep, state),
					PrevSteps:   StepStack{},
				}
			}

			b.processUserSession(msg)
		}
	}()

	go func() {
		for msg := range b.sendMsgChan {
			if err := b.clients[msg.Source].SendMessage(msg); err != nil {
				fmt.Println(err)
			}
		}
	}()
}

func (b *Bot) Shutdown() {
	for _, client := range b.clients {
		client.Disconnect()
	}
}

func (b *Bot) createStep(step StepType, state *entities.UserState) Step {
	switch step {
	case MainMenuStep:
		return &MainMenu{
			StepBase: StepBase{logger: b.logger, State: state},
		}
	case MainMenuRequestStep:
		return &YesNo{
			StepBase: StepBase{logger: b.logger},
			question: Question{Text: "Вурнуться в главное меню?"}, yesStep: MainMenuStep, noStep: EmptyStep,
		}
	case CitySelectionStep:
		return &CitySelection{
			StepBase:     StepBase{logger: b.logger, State: state, DbAdapter: b.dbAdapter},
			checkService: true,
			filter:       true,
			nextStep:     MasterSelectionStep,
			errStep:      EmptyStep,
		}
	case CityPromptStep:
		return &CityPrompt{
			StepBase: StepBase{logger: b.logger, State: state, DbAdapter: b.dbAdapter},
		}
	case ServiceCategorySelectionStep:
		return &ServiceCategorySelection{
			StepBase: StepBase{logger: b.logger, State: state, DbAdapter: b.dbAdapter},
			filter:   true,
			errStep:  EmptyStep,
		}
	case ServiceSelectionStep:
		return &ServiceSelection{
			StepBase: StepBase{logger: b.logger, State: state, DbAdapter: b.dbAdapter},
		}
	case MasterSelectionStep:
		return &MasterSelection{StepBase: StepBase{logger: b.logger, State: state, DbAdapter: b.dbAdapter}}
	case FinalStep:
		return &Final{StepBase{logger: b.logger, State: state, DbAdapter: b.dbAdapter}}
	case MasterStep:
		return &YesNo{
			StepBase: StepBase{logger: b.logger, State: state, DbAdapter: b.dbAdapter},
			question: Question{Text: "Хотите зарегистрироваться как мастер?"},
			yesStep:  RegistrationStep,
			noStep:   MainMenuStep,
		}
	case RegistrationStep:
		return &Prompt{
			StepBase: StepBase{logger: b.logger, State: state, DbAdapter: b.dbAdapter},
			question: Question{Text: "Как вас называть?", Field: "name"},
			nextStep: RegistrationStepCity,
			errStep:  RegistrationStep,
		}
	case RegistrationStepCity:
		return &CitySelection{
			StepBase: StepBase{logger: b.logger, State: state, DbAdapter: b.dbAdapter},
			nextStep: RegistrationStepService,
			errStep:  EmptyStep,
		}
	case RegistrationStepService:
		return &ServiceSelection{
			StepBase: StepBase{logger: b.logger, State: state, DbAdapter: b.dbAdapter},
		}
	case RegistrationFinalStep:
		return &RegistrationFinal{StepBase: StepBase{logger: b.logger, State: state, DbAdapter: b.dbAdapter}}
	case EmptyStep:
		return nil
	case AdminStep:
		return &Admin{StepBase: StepBase{logger: b.logger, State: state, DbAdapter: b.dbAdapter}}
	case AdminServiceCategorySelectionStep:
		return &ServiceCategorySelection{
			StepBase:   StepBase{logger: b.logger, State: state, DbAdapter: b.dbAdapter},
			addSrvMode: true,
			filter:     false,
			errStep:    EmptyStep,
		}
	case AddServiceCategoryStep:
		return &AddServiceCategory{StepBase: StepBase{logger: b.logger, State: state, DbAdapter: b.dbAdapter}}
	case AddServiceStep:
		return &AddService{StepBase: StepBase{logger: b.logger, State: state, DbAdapter: b.dbAdapter}}
	case AddCityStep:
		return &AddCity{StepBase: StepBase{logger: b.logger, State: state, DbAdapter: b.dbAdapter}}
	case TestStep:
		return &Test{StepBase: StepBase{logger: b.logger, State: state, DbAdapter: b.dbAdapter}}
	default:
		return &MainMenu{StepBase: StepBase{logger: b.logger, State: state}}
	}
}

func (b *Bot) send(msg *ma.Message) bool {
	if msg == nil {
		return false
	}
	b.sendMsgChan <- msg
	b.logger.Infof("sending a message: %s", msg.Text)
	return true
}

func (b *Bot) processMessage(msg *ma.Message) {
	curStep := b.userSessions[msg.UserID].CurrentStep
	state := b.userSessions[msg.UserID].State
	if !curStep.IsInProgress() {
		b.send(curStep.Request(msg))
	} else {
		res, next := curStep.ProcessResponse(msg)
		b.send(res)

		switch step := b.createStep(next, state); next {
		case PreviousStep:
			var prevStep Step
			if b.userSessions[msg.UserID].PrevSteps.Empty() {
				prevStep = b.createStep(MainMenuStep, state)
			} else {
				prevStep = b.userSessions[msg.UserID].PrevSteps.Top()
				b.userSessions[msg.UserID].PrevSteps.Pop()
			}
			prevStep.Reset()

			b.send(prevStep.Request(msg))
			b.userSessions[msg.UserID].CurrentStep = prevStep
		case EmptyStep:
		case MainMenuStep:
			b.send(step.Request(msg))
			b.userSessions[msg.UserID].CurrentStep = step
			b.userSessions[msg.UserID].PrevSteps.Clear()
		case MainMenuRequestStep:
			time.Sleep(1 * time.Second)
			fallthrough
		default:
			b.send(step.Request(msg))

			if step.IsCallBackStep() {
				curStep.SetInProgress(true)
				b.userSessions[msg.UserID].СurrentCallbackStep = step
			} else {
				b.userSessions[msg.UserID].CurrentStep = step
			}

			if curStep.IsCallBackStep() {
				b.userSessions[msg.UserID].PrevCallBackSteps.Push(curStep)
			} else {
				b.userSessions[msg.UserID].PrevSteps.Push(curStep)
			}
		}
	}
}

func (b *Bot) processCallback(msg *ma.Message) {
	curStep := b.userSessions[msg.UserID].СurrentCallbackStep
	state := b.userSessions[msg.UserID].State
	if !curStep.IsInProgress() {
		b.send(curStep.Request(msg))
	} else {
		res, next := curStep.ProcessResponse(msg)
		b.send(res)

		switch step := b.createStep(next, state); next {
		case PreviousStep:
			if b.userSessions[msg.UserID].PrevCallBackSteps.Empty() {
				return
			}

			prevStep := b.userSessions[msg.UserID].PrevCallBackSteps.Top()
			b.userSessions[msg.UserID].PrevCallBackSteps.Pop()
			prevStep.Reset()

			b.send(prevStep.Request(msg))
			b.userSessions[msg.UserID].СurrentCallbackStep = prevStep
		case EmptyStep:
		case MainMenuRequestStep:
			time.Sleep(1 * time.Second)
			fallthrough
		default:
			b.send(step.Request(msg))

			if step.IsCallBackStep() {
				b.userSessions[msg.UserID].СurrentCallbackStep = step
			} else {
				b.userSessions[msg.UserID].CurrentStep = step
			}

			if curStep.IsCallBackStep() {
				b.userSessions[msg.UserID].PrevCallBackSteps.Push(curStep)
			} else {
				b.userSessions[msg.UserID].PrevSteps.Push(curStep)
			}
		}
	}
}

func (b *Bot) processUserSession(msg *ma.Message) {

	if msg.Type == ma.CALLBACK {
		b.processCallback(msg)
	} else {
		b.processMessage(msg)
	}
}
