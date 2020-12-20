package infra

var apiInitializerRegister *InitializeRegister = new(InitializeRegister)

type WebApiStarter struct {
	BaseStarter
}

//注册WEB API初始化对象

func RegisterApi(ai Initializer) {
	apiInitializerRegister.Register(ai)
}

//获取注册的web api初始化对象
func GetApiInitializers() []Initializer {
	return apiInitializerRegister.Initializers
}

func (w *WebApiStarter) Setup(ctx StarterContext) {
	for _, v := range GetApiInitializers() {
		v.Init()
	}
}
