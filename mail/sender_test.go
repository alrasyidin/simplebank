package mail

import (
	"testing"

	"github.com/alrasyidin/simplebank-go/util"
	"github.com/stretchr/testify/require"
)

func TestSendEmail(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	config, err := util.LoadConfig("../")
	require.NoError(t, err)

	gm := NewGmailSender(config.GmailName, config.GmailUsername, config.GmailPassword)

	subject := "Test email"
	content := `
    <h1>Hello World</h1>
    <p>
      Hari ini hari senin, semangat kerjanya
    </p>
  `
	to := []string{"hafidhpradiptaarrasyid@gmail.com"}
	attachFiles := []string{"../README.md"}
	err = gm.SendEmail(subject, content, to, nil, nil, attachFiles)

	require.NoError(t, err)
}
