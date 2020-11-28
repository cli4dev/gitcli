package update

import (
	"fmt"
	"net/smtp"
	"net/textproto"
	"strings"

	"github.com/micro-plat/cli/cmds"
	"github.com/micro-plat/cli/logs"
	"github.com/micro-plat/lib4go/types"
	"github.com/urfave/cli"

	"github.com/jordan-wright/email"
)

func init() {
	cmds.Register(
		cli.Command{
			Name:   "email",
			Usage:  "发送电子邮件",
			Action: upload,
			Flags: []cli.Flag{
				cli.StringFlag{
					Name:     "from,f",
					Required: true,
					Usage:    "来源",
				},
				cli.StringFlag{
					Name:     "to,t",
					Required: true,
					Usage:    "目标",
				},
				cli.StringFlag{
					Name:  "subject,s",
					Usage: "主题",
				},
				cli.StringFlag{
					Name:  "text",
					Usage: "内容",
				},
				cli.StringFlag{
					Name:  "att,a",
					Usage: "附件",
				},
			},
		})
}

//pull 根据传入的路径(分组/仓库)拉取所有项目
func upload(c *cli.Context) (err error) {
	from := c.String("from")
	to := c.String("to")
	subject := c.String("subject")
	text := c.String("text")
	att := c.String("att")

	f := strings.Split(from, ",")

	if len(f) != 2 || f[0] == "" || f[1] == "" {
		return fmt.Errorf("邮件发送人密码不能为空，格式:sender@mail.com,pwd")
	}
	if len(f) != 2 || f[0] == "" || f[1] == "" || to == "" {
		return fmt.Errorf("参数不能为空from:%s to:%s", from, to)
	}
	mail := &email.Email{
		To:      []string{to},
		From:    f[0],
		Subject: types.GetString(subject, "无主题"),
		Text:    []byte(text),
		Headers: textproto.MIMEHeader{},
	}
	if att != "" {
		mail.AttachFile(att)
	}
	logs.Log.Info("正在发送...")

	err = mail.Send("smtp.exmail.qq.com:587", smtp.PlainAuth("", f[0], f[1], "smtp.exmail.qq.com"))
	if err != nil {
		return fmt.Errorf("发送失败:%w", err)
	}
	logs.Log.Info("邮件发送成功")
	return nil

}
