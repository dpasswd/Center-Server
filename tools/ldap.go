package tools

import (
	"fmt"
	"time"

	"github.com/go-ldap/ldap"
)

type LDAPConfig struct {
	Addr         string         //ldap 地址:端口
	BindUserName string         //用户名
	BindPassword string         //密码
	SearchDN     string         //用户OU
	Attributes   LdapAttributes //显示字段
}

type LdapAttributes struct {
	UsernameKey string `json:"username"` //ldap用户名的key
	NicknameKey string `json:"nickname"` //ldap昵称的key
	EmailKey    string `json:"email"`    //ldap email的key
	PhoneKey    string `json:"phone"`    //ldap电话的key
}

type UserInfo struct {
	Username string
	Nickname string
	Email    string
	Phone    string
}

type LDAPService struct {
	Conn   *ldap.Conn
	Config LDAPConfig
}

// 创建ldap连接
func NewLDAPService(config LDAPConfig) (*LDAPService, error) {
	conn, err := ldap.Dial("tcp", config.Addr)
	if err != nil {
		return nil, err
	}
	conn.SetTimeout(5 * time.Second)
	// NOTE(chenjun): 暂时先不skip verify
	// err = conn.StartTLS(&tls.Config{InsecureSkipVerify: true})
	// if err != nil {
	//  return nil, err
	// }
	err = conn.Bind(config.BindUserName, config.BindPassword)
	if err != nil {
		return nil, err
	}
	return &LDAPService{Conn: conn, Config: config}, nil
}

// Login 登录并获取用户信息
func (l *LDAPService) Login(userName, password string) (bool, UserInfo, error) {
	var userInfo UserInfo
	defer l.Conn.Close()
	searchRequest := ldap.NewSearchRequest(
		l.Config.SearchDN,
		ldap.ScopeWholeSubtree, ldap.NeverDerefAliases, 0, 0, false,
		fmt.Sprintf("(&(objectClass=*)(%s=%s))", l.Config.Attributes.UsernameKey, userName),
		[]string{"dn",
			l.Config.Attributes.EmailKey,
			l.Config.Attributes.NicknameKey,
			l.Config.Attributes.UsernameKey,
			l.Config.Attributes.PhoneKey},
		nil,
	)
	sr, err := l.Conn.Search(searchRequest)
	if err != nil {
		return false, userInfo, err
	}
	if len(sr.Entries) != 1 {
		return false, userInfo, fmt.Errorf("User does not exist or too many entries returned")
	}
	userDN := sr.Entries[0].DN
	err = l.Conn.Bind(userDN, password)
	if err != nil {
		return false, userInfo, err
	}
	err = l.Conn.Bind(l.Config.BindUserName, l.Config.BindPassword)
	if err != nil {
		return false, userInfo, nil
	}
	for _, v := range sr.Entries[0].Attributes {
		switch v.Name {
		case l.Config.Attributes.EmailKey:
			userInfo.Email = v.Values[0]
		case l.Config.Attributes.NicknameKey:
			userInfo.Nickname = v.Values[0]
		case l.Config.Attributes.UsernameKey:
			userInfo.Username = v.Values[0]
		case l.Config.Attributes.PhoneKey:
			userInfo.Phone = v.Values[0]
		}
	}
	return true, userInfo, nil
}
