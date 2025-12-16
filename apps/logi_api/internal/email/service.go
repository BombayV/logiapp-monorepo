package email

import (
	"fmt"

	"github.com/resend/resend-go/v2"
)

type Service struct {
	client *resend.Client
}

func NewService(apiKey string) *Service {
	client := resend.NewClient(apiKey)
	return &Service{client: client}
}

func (s *Service) SendSurveyEmail(to []string, surveyLink string) error {
	htmlContent := fmt.Sprintf(`
		<p>Hola,</p>
		<p>Gracias por confiar en nosotros. Por favor ayúdanos a mejorar respondiendo esta breve encuesta sobre tu pedido:</p>
		<p><a href="%s">Responder Encuesta</a></p>
		<p>Si el enlace no funciona, copia y pega esta dirección en tu navegador:</p>
		<p>%s</p>
	`, surveyLink, surveyLink)

	params := &resend.SendEmailRequest{
		From:    "onboarding@resend.dev",
		To:      to,
		Subject: "Encuesta de Satisfacción - LogiApp",
		Html:    htmlContent,
	}

	_, err := s.client.Emails.Send(params)
	return err
}
