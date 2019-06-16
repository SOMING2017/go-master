package message

import (
	"fmt"
	"regexp"
)

//规定：
//不为空，且全部是数字或字母
func VerifyChannelID(channelID string) error {
	if channelID == "" {
		return fmt.Errorf("账户不能为空")
	}
	maBool, err := regexp.MatchString("^\\w+$", channelID)
	if err != nil {
		return err
	}
	if !maBool {
		return fmt.Errorf("通道ID不符合规则")
	}
	return nil
}
