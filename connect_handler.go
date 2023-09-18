package hcpvaultengine

type AuthStatus struct {
	Token string //should this be encrypted? how?
}

var (
	_ ConnectHandler = (*interactiveConnectHandler)(nil)
	_ ConnectHandler = (*nonInteractiveConnectHandler)(nil)
)

type ConnectHandler interface {
	// should we parse this before [during the Run function common flag parsin] into a struct that's specific to the different handlers?
	Connect(args []string) (AuthStatus, error)
}

type interactiveConnectHandler struct{}

func (h *interactiveConnectHandler) Connect(args []string) (AuthStatus, error) {
	return AuthStatus{}, nil
}

type nonInteractiveConnectHandler struct{}

func (h *nonInteractiveConnectHandler) Connect(args []string) (AuthStatus, error) {
	return AuthStatus{}, nil
}
