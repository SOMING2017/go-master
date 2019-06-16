package goChan

type T = []byte

//忽略通道再次关闭异常
func SafeClose(ch chan T) (justClosed bool) {
	defer func() {
		if recover() != nil {
			justClosed = false
		}
	}()
	close(ch)
	return true
}

//忽略向已关闭通道发送数据异常
func SafeSend(ch chan T, value T) (closed bool) {
	defer func() {
		if recover() != nil {
			closed = true
		}
	}()
	ch <- value
	return false
}

// ---------------------
// 作者：dengming0922
// 来源：CSDN
// 原文：https://blog.csdn.net/dengming0922/article/details/80904235
// 版权声明：本文为博主原创文章，转载请附上博文链接！

func SafeSendMessage(ch chan T, value T) {
	select {
	case ch <- value:
	default:
		//no deal
	}
}
