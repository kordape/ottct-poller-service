package event

type FakeNews struct {
	EntityId  string
	Timestamp string
	Content   string
}

type SendFakeNewsEventFn func(events []FakeNews) error

func SendFakeNewsEventFnBuilder() SendFakeNewsEventFn {
	return func(events []FakeNews) error {
		// TODO: replace with sending events to SQS
		return nil
	}
}
