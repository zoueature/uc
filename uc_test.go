package ucgo

import (
	"github.com/golang/mock/gomock"
	"github.com/jiebutech/uc/cache"
	"github.com/jiebutech/uc/model"
	"github.com/jiebutech/uc/sender"
	"github.com/jiebutech/uc/types"
	"testing"
)

func TestMain(m *testing.M) {
	m.Run()
}

func password() Password {
	return Password{
		Salt:     "123213132131",
		Password: "123456",
	}
}
func TestNewInternalClient(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockUser := model.NewMockUserEntity(ctrl)
	mockUser.EXPECT().SetUsername("username").AnyTimes()
	mockUser.EXPECT().GetUserName().AnyTimes().Return("username")
	mockUser.EXPECT().GetID().AnyTimes().Return(int64(123))
	mockUser.EXPECT().SetLoginType(types.EmailLogin).AnyTimes()
	mockUser.EXPECT().GetLoginType().AnyTimes()
	mockUser.EXPECT().SetIdentify("767955912@qq.com").AnyTimes().Return()
	mockUser.EXPECT().GetIdentify().AnyTimes().Return("767955912@qq.com")
	mockUser.EXPECT().GetPassword().AnyTimes().Return(password())
	mockUser.EXPECT().ToMap().AnyTimes().Return(map[string]interface{}{
		"id": 123,
	})

	mockUserRepo := model.NewMockUserResource(ctrl)
	mockUserRepo.EXPECT().GenUser().AnyTimes().Return(mockUser)
	mockUserRepo.EXPECT().GetUserByIdentify(mockUser).AnyTimes().Return(nil)
	mockUserRepo.EXPECT().GetUserByUsername(mockUser).AnyTimes().Return(nil)

	cli := NewUserClient(
		cache.NewRedisCache("192.168.1.202:6379", ""),
		sender.NewEmailSender(&sender.EmailConfig{
			From:     "zoueature@gmail.com",
			Host:     "smtp.gmail.com",
			Port:     465,
			Username: "zoueature",
			Password: "dvyjaualoktcrsxx",
		}, sender.DefaultMailMessage()),
		mockUserRepo,
		DefaultJwtEncoder("12321313", 32132131),
	)
	err := cli.SendSmsCode(types.RegisterCodeType, UserIdentify{
		App:      "",
		Type:     "",
		Identify: "",
	})
	if err != nil {
		t.Fatal(err)
	}

	//token, info, err := cli.Login(types.EmailLogin, "767955912@qq.com", "123456")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//token, info, err = cli.LoginByUsername(types.EmailLogin, "username", "123456")
	//if err != nil {
	//	t.Fatal(err)
	//}
	//t.Log(token)
	//t.Log(info.ToMap())
}
