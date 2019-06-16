package user

import (
	"../im/channel/main"
)

func bindingUserWorldChannel(userID string) error {
	return channel.BindingUserWorldChannel(userID)
}
