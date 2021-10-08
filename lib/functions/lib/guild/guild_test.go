package guild

import (
	"testing"
)

func TestValidate(t *testing.T) {

	t.Run("succsess if all fields are supported", func(t *testing.T) {
		if err := in().validate(); err != nil {
			t.Error(err)
		}
	})
	t.Run("rinkeby should be supported", func(t *testing.T) {
		in := in()
		in.Network = "rinkeby"
		if err := in.validate(); err != nil {
			t.Error(err)
		}
	})
	t.Run("polygon should not be supported", func(t *testing.T) {
		in := in()
		in.Network = "polygon"
		if err := in.validate(); err == nil {
			t.Error("unexpected")
		}
	})
	t.Run("error if roleid is empty", func(t *testing.T) {
		in := in()
		in.RoleID = ""
		if err := in.validate(); err == nil {
			t.Error("unexpected")
		}
	})
}

func in() ContractInfo {
	return ContractInfo{
		RoleID:          "test",
		ContractAddress: "test",
		Network:         "mainnet",
		GuildID:         "test",
	}
}

//func TestNotify(t *testing.T) {
//	svc, err := New(context.Background(), ssm.New(), nil)
//	if err != nil {
//		t.Error(err)
//	}
//	err = svc.notify(context.Background(), "892441777808765052", ContractInfo{
//		RoleID:          "test",
//		ContractAddress: "test",
//		Network:         "rinkeby",
//		GuildID:         "test",
//	})
//	if err != nil {
//		t.Error(err)
//	}
//}
//
//func TestError(t *testing.T) {
//	svc, err := New(context.Background(), ssm.New(), nil)
//	if err != nil {
//		t.Error(err)
//	}
//	err = svc.error(context.Background(), "892441777808765052", errors.New("error message"))
//	if err != nil {
//		t.Error(err)
//	}
//}
