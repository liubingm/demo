package util

func BuildLogAnalysisPrompt(logs []string) (string, []Message) {
	var messages []Message

	// System prompt preparation
	systemPrompt :=
		`
		Human: You are a AI trained MQTT Operations Expert, please follow below instruction and help me to anlyze my mqtt log.
			1. please help analysis the error reason if there are erros.
			2. locate the error log lines.
			3. according to the reason, list several solutions.

		Please out put json format as below:\n
		{
			"Reason": "",
			"Solution": ""
		}
		`
	var userPrompt string

	userPrompt = "Here is some AWS IoTCore log data, please help me to analyze them"

	for _, log := range logs {

		userPrompt = userPrompt + log + "/n/n"
	}

	assistantPrompt := `{
		"Reason":`

	messages = []Message{
		{Role: "user", Content: userPrompt},
		{Role: "assistant", Content: assistantPrompt}}

	return systemPrompt, messages
}
