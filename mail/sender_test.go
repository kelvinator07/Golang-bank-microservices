package mail

import (
	"testing"

	"github.com/kelvinator07/golang-bank-microservices/util"
	"github.com/stretchr/testify/assert"
)

func TestSendEmailWithGmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}

	config, err := util.LoadConfig("./../app.env")
	assert.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "A test email from golang bank"
	content := `
	<h1>Hello World</h1>
	<p>This is a test message from <a href="https://github.com/kelvinator07">Geeky Kel Github</a></p>
	`
	to := []string{"isievwore.kelvin@gmail.com"}
	attachFiles := []string{"../comments.txt"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	assert.NoError(t, err)
}
