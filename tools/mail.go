package tools

import (
	"crypto/tls"
	"encoding/base64"
	"errors"
	"log"
	"mime"
	"net"
	"net/mail"
	"net/smtp"
	"strings"
)

// 错误列表
var (
	//ErrNoMailFrom address of from user is empty
	ErrNoMailFrom = errors.New("from field is empty")
	//ErrNoMailTo address of to is empty
	ErrNoMailTo = errors.New("to field is empty")
)

// 邮件参数
type EmailParam struct {
	SmtpAddr    string `json:"smtpAddr"`
	SmtpSec     string `json:"smtpSec"`
	SmtpPort    string `json:"smtpPort"`
	SmtpAccount string `json:"smtpAccount"`
	SmtpPwd     string `json:"smtpPwd"`
	Toers       string `json:"toers"`
}

// 格式化邮件地址
func parseAddress(a []string) ([]*mail.Address, error) {
	addrs := make([]*mail.Address, len(a))
	for i := 0; i < len(a); i++ {
		addr, err := mail.ParseAddress(a[i])
		if err != nil {
			return nil, err
		}
		addrs[i] = addr
	}
	return addrs, nil
}

// 发送注册激活邮件
func (e *EmailParam) SendRegisterMail(toMail string, url string) (err error) {
	var (
		sub = "hdPassword 账号激活"
		msg = []byte(`<style>
		html{-ms-text-size-adjust:100%;-webkit-text-size-adjust:100%}body{line-height:1.6;font-family:"Helvetica Neue",Helvetica,Arial,sans-serif;font-size:16px}body,dd,dl,fieldset,h1,h2,h3,h4,h5,ol,p,textarea,ul{margin:0}button,fieldset,input,legend,textarea{padding:0}button,input,select,textarea{font-family:inherit;font-size:100%;margin:0}ol,ul{padding-left:0;list-style-type:none}a img,fieldset{border:0}a{text-decoration:none}.radius_avatar{display:inline-block;background-color:#FFF;padding:3px;border-radius:50%;-moz-border-radius:50%;-webkit-border-radius:50%;overflow:hidden;vertical-align:middle}.radius_avatar img{display:block;width:100%;height:100%;border-radius:50%;-moz-border-radius:50%;-webkit-border-radius:50%;background-color:#EEE}.btn_app{margin-top:10px;position:relative;display:block;margin-left:auto;margin-right:auto;padding-left:14px;padding-right:14px;-webkit-box-sizing:border-box;box-sizing:border-box;font-size:16px;text-align:center;text-decoration:none;color:#FFF;line-height:2.625;border-radius:5px;-webkit-tap-highlight-color:transparent;overflow:hidden}.btn_app:after{content:" ";width:200%;height:200%;position:absolute;top:0;left:0;border:1px solid rgba(0,0,0,.2);-webkit-transform:scale(.5);transform:scale(.5);-webkit-transform-origin:0 0;transform-origin:0 0;-webkit-box-sizing:border-box;box-sizing:border-box;border-radius:10px}.btn_app_primary{background-color:#42C642}.btn_app_primary:link,.btn_app_primary:visited{color:#FFF}.btn_app_primary:active{color:rgba(255,255,255,.6)}.btn_app_default{background-color:#F7F7F7;color:#454545}.btn_app_default:link,.btn_app_default:visited{color:#454545}.btn_app_default:active{color:#C9C9C9}.skin_app_default{background-image:url();-webkit-background-size:100% auto;background-size:100% auto;background-position:50% 0;background-repeat:no-repeat;background-color:#FFF}body,html{position:relative;height:100%}a,a:link,a:visited{color:#42C642}.mail_area{text-align:center;height:100%;-webkit-box-sizing:border-box;box-sizing:border-box;display:-webkit-box;display:-webkit-flex;display:-ms-flexbox;display:flex;-webkit-box-align:center;-webkit-align-items:center;-ms-flex-align:center;align-items:center;-webkit-box-pack:center;-webkit-justify-content:center;-ms-flex-pack:center;justify-content:center;font-family:"Helvetica Neue","Hiragino Sans GB","Microsoft YaHei","\9ED1\4F53",Arial,sans-serif}.mail{position:relative;display:inline-block;width:80%;margin-top:-150px;text-align:left}.mail_pc{background-color:#E6E6EA;display:block}.mail_pc .mail{width:850px;margin:45px 0;box-shadow:0 0 25px 5px rgba(0,0,0,.09);-moz-box-shadow:0 0 25px 5px rgba(0,0,0,.09);-webkit-box-shadow:0 0 25px 5px rgba(0,0,0,.09);background-color:#FFF;border-radius:8px;-moz-border-radius:8px;-webkit-border-radius:8px;overflow:hidden}.mail_pc .mail_inner{padding:17% 16%}.mail_pc .mail_msg .btn_app{width:225px}.pic_skin_top{position:absolute;top:0;left:0;width:145px;height:175px;background:url() no-repeat}.pic_skin_bottom{position:absolute;bottom:0;right:0;width:300px;height:265px;background:url() no-repeat}h1{font-weight:400;position:absolute;right:48px;top:48px;line-height:300px;overflow:hidden;width:314px;height:32px;background:url() no-repeat}.mail_info{padding:1.7em 0 0 56px;margin-top:4.3em;position:relative;border-top:1px #BBBBBD dashed;font-size:14px}.mail_info .radius_avatar{width:38px;height:38px;padding:0;position:absolute;top:1.7em;left:0}.mail_info strong{font-weight:400}.mail_info p{color:#C1C1C3;margin-top:-.34em;font-size:12px}.mail_msg h2{font-weight:400;font-size:20px;color:#1D1D26;padding:1.34em 0 .6em}.mail_msg p{margin-bottom:24px}.mail_msg .btn_app{margin-top:45px}#app_mail .mail_msg .btn_app,#app_mail .mail_msg .btn_app:link,#app_mail .mail_msg .btn_app:visited{text-decoration:none}
		</style>
		<div class="mail_area mail_pc" id="app_mail">
			<div class="mail">
				<div class="mail_inner">
					<h1>test</h1>
					<div class="mail_msg">
						<p>
							<br> 你好: </br>
							<br> 感谢你注册HD Password </br>
							<br> 你的登录邮箱为：` + toMail + `。请点击以下链接激活账号， </br>
						</p>
						<p>
							<a href="` + url + `" target="_blank">` + url + `</a>
						</p>
						<div class="mail_info", align=right>
							<strong>电魂运维团队</strong>
						</div>
					</div>
					<div class="pic_skin_top"></div>
					<div class="pic_skin_bottom"></div>
				</div>
			</div>
		</div>`)
	)
	var (
		addr     = net.JoinHostPort(e.SmtpAddr, e.SmtpPort)
		from     = e.SmtpAccount
		password = e.SmtpPwd
		to       = []string{toMail}
		cc       = []string{}
		// cc       = []string{"7@qq.com"}
		bcc = []string{}
	)
	if e.SmtpSec != "NULL" {
		go func() {
			err = Sendmail(addr, from, password, to, cc, bcc, sub, msg)
			log.Println(err)
		}()
	} else {
		go func() {
			err = SkipVerifyTLS(addr, from, password, to, cc, bcc, sub, msg)
			log.Println(err)
		}()
	}
	return err
}

// 格式化 message
func newMessage(from *mail.Address, to, cc, bcc []*mail.Address, sub string, msg []byte) (string, error) {
	var buf strings.Builder
	// write subject
	buf.WriteString("Subject: ")
	buf.WriteString(mime.BEncoding.Encode("utf-8", sub))
	buf.WriteString("\r\n")
	//write mail from
	buf.WriteString("From: ")
	buf.WriteString(from.String())
	buf.WriteString("\r\n")
	// write rcpt to
	buf.WriteString("To: ")
	for i := 0; i < len(to); i++ {
		buf.WriteString(to[i].String())
		if i != len(to)-1 {
			buf.WriteByte(',')
		}
	}
	buf.WriteString("\r\n")
	// write cc
	buf.WriteString("Cc: ")
	for i := 0; i < len(cc); i++ {
		buf.WriteString(cc[i].String())
		if i != len(to)-1 {
			buf.WriteByte(',')
		}
	}
	buf.WriteString("\r\n")
	// write Bcc
	buf.WriteString("Bcc: ")
	for i := 0; i < len(bcc); i++ {
		buf.WriteString(bcc[i].String())
		if i != len(to)-1 {
			buf.WriteByte(',')
		}
	}
	buf.WriteString("\r\n")
	// write content-type
	buf.WriteString("Content-Type: text/html; charset=UTF-8\r\n")
	buf.WriteString("Content-Transfer-Encoding: base64\r\n")
	buf.WriteString("\r\n")
	buf.WriteString(base64.StdEncoding.EncodeToString(msg))
	buf.WriteString("\r\n")
	return buf.String(), nil
}

//Sendmail use smtp.SendMail
func Sendmail(addr, from, password string, to, cc, bcc []string, sub string, msg []byte) error {
	//check addr with net.SplitHostPort
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}

	if from == "" {
		return ErrNoMailFrom
	}

	if len(to) < 1 {
		return ErrNoMailTo
	}
	mailfrom, err := mail.ParseAddress(from)
	if err != nil {
		return err
	}
	rcptto, err := parseAddress(to)
	if err != nil {
		return err
	}
	ccto, err := parseAddress(cc)
	if err != nil {
		return err
	}
	bccto, err := parseAddress(bcc)
	if err != nil {
		return err
	}

	content, err := newMessage(mailfrom, rcptto, ccto, bccto, sub, msg)
	if err != nil {
		return err
	}
	rcptto = append(rcptto, ccto...)
	rcptto = append(rcptto, bccto...)
	mailto := make([]string, len(rcptto))
	for i := 0; i < len(rcptto); i++ {
		mailto[i] = rcptto[i].Address
	}
	// fmt.Println(mailfrom.Address)
	// fmt.Println(password)

	if password == "" {
		if err := smtp.SendMail(addr, nil, mailfrom.Address, mailto, []byte(content)); err != nil {
			return err
		}
	} else {
		//PlainAuth
		a := smtp.PlainAuth("", mailfrom.Address, password, host)
		// a := smtp.CRAMMD5Auth(mailfrom.Address, password)
		if err := smtp.SendMail(addr, a, mailfrom.Address, mailto, []byte(content)); err != nil {
			return err
		}
	}
	return nil
}

//SkipVerifyTLS sendmail skip TLS verify, it only should be used on localhost
func SkipVerifyTLS(addr, from, password string, to, cc, bcc []string, sub string, msg []byte) error {
	//check addr with net.SplitHostPort
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return err
	}

	if from == "" {
		return ErrNoMailFrom
	}

	if len(to) < 1 {
		return ErrNoMailTo
	}

	c, err := smtp.Dial(addr)
	if err != nil {
		return err
	}
	defer c.Close()
	//skip tls verify
	config := &tls.Config{
		InsecureSkipVerify: true,
		ServerName:         host,
	}
	if err = c.StartTLS(config); err != nil {
		return err
	}

	mailfrom, err := mail.ParseAddress(from)
	if err != nil {
		return err
	}
	//auth password
	if password != "" {
		a := smtp.PlainAuth("", mailfrom.Address, password, host)
		if err := c.Auth(a); err != nil {
			return err
		}
	}
	if err = c.Mail(mailfrom.Address); err != nil {
		return err
	}

	rcptto, err := parseAddress(to)
	if err != nil {
		return err
	}
	ccto, err := parseAddress(cc)
	if err != nil {
		return err
	}
	bccto, err := parseAddress(bcc)
	if err != nil {
		return err
	}

	for _, rcpt := range rcptto {
		if err := c.Rcpt(rcpt.Address); err != nil {
			return err
		}
	}
	for _, rcpt := range ccto {
		if err := c.Rcpt(rcpt.Address); err != nil {
			return err
		}
	}
	for _, rcpt := range bccto {
		if err := c.Rcpt(rcpt.Address); err != nil {
			return err
		}
	}

	content, err := newMessage(mailfrom, rcptto, ccto, bccto, sub, msg)
	if err != nil {
		return err
	}
	//wc is io.WriteCloser
	wc, err := c.Data()
	if err != nil {
		return err
	}

	// write email content to wc
	if _, err := wc.Write([]byte(content)); err != nil {
		return err
	}
	if err = wc.Close(); err != nil {
		return err
	}
	return c.Quit()
}
