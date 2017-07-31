package email

import (
        "github.com/aws/aws-sdk-go/aws/session"
        "github.com/aws/aws-sdk-go/service/ses"
        "github.com/aws/aws-sdk-go/aws"
        "github.com/op/go-logging"
        "html/template"
)

type EmailService interface {
        SendEmail(mail string, content string, subject string) error
        Templates() *template.Template
}

type DefaultService struct {
        sess *session.Session
        replyAddr string
}

const encoding = "UTF-8"

func (df DefaultService) Templates() *template.Template{
        return template.Must(template.New("confirm_account.html").ParseFiles("email/templates/confirm_account.html", "email/templates/reset_password.html"))
}

func (df DefaultService) SendEmail(email string, content string, subject string) error {

        log := logging.MustGetLogger("[email]")
        svc := ses.New(df.sess)

        params := &ses.SendEmailInput{
                Destination: &ses.Destination{
                        ToAddresses: []*string{
                                aws.String(email),
                        },
                },
                Message: &ses.Message{
                        Body: &ses.Body{
                                Html: &ses.Content{
                                        Data:    aws.String(content),
                                        Charset: aws.String(encoding),
                                },
                        },
                        Subject: &ses.Content{
                                Data:    aws.String(subject),
                                Charset: aws.String(encoding),
                        },
                },
                Source: aws.String(df.replyAddr),
        }

        _, err := svc.SendEmail(params)

        if err != nil {
                // Print the error, cast err to awserr.Error to get the Code and
                log.Error("Error while sending email " + err.Error())
                return err
        }

        return nil
}

func Create(awsEmailRegion string, replyAddr string) (EmailService, error) {

        sess, err := session.NewSession(&aws.Config{Region: aws.String(awsEmailRegion)})
        if err != nil {
                return nil, err
        }

        return DefaultService{sess:sess, replyAddr: replyAddr}, nil
}




