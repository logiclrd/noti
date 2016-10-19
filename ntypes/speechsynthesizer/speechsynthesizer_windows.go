package speechsynthesizer

type Notification struct {
	Text string
	// Rate is from -10 to 10. -10 is slowest.
	Rate  int
	Voice string
}

// $n = New-Object System.Speech.Synthesis.SpeechSynthesizer
// $n.Speak('Hello')
