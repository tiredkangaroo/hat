package database

const (
	// TriggerIncomingRequest is for incoming requests to the server.
	TriggerIncomingRequest = "incoming_request"
	// TriggerRecievedMITMRequest is called when the request is secure but MITM is able to
	// capture the request intended for the host.
	TriggerRecievedMITMRequest = "mitm_handled"
)
