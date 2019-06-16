package user

import (
	"fmt"
	"regexp"
)

//规定：
//6-9位数字字母，并且第一位必须为字母
func VerifyUser(user string) error {
	if user == "" {
		return fmt.Errorf("账户不能为空")
	}
	maBool, err := regexp.MatchString("^[A-Za-z]\\w{5,8}$", user)
	if err != nil {
		return err
	}
	if !maBool {
		return fmt.Errorf("账户不符合规则")
	}
	return nil
}

//规定：
//8-16位数字字母
func VerifyPassword(password string) error {
	if password == "" {
		return fmt.Errorf("密码不能为空")
	}
	maBool, err := regexp.MatchString("^\\w{8,16}$", password)
	if err != nil {
		return err
	}
	if !maBool {
		return fmt.Errorf("密码不符合规则")
	}
	return nil
}

//规定：
//2-6位汉字
func VerifyName(name string) error {
	if name == "" {
		return fmt.Errorf("用户名不能为空")
	}
	maBool, err := regexp.MatchString("^[\\p{Han}]{2,6}$", name)
	if err != nil {
		return err
	}
	if !maBool {
		return fmt.Errorf("用户名不符合规则")
	}
	return nil
}
