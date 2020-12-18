package infra

import "github.com/tietang/props/kvs"

const (
	KeyProps = "_conf"
)

//资源启动器上下文
//用来在服务资源初始化、安装、启动和停止的生命周期中变量和对象的传递
type StarterContext map[string]interface{}

func (s StarterContext) Props() kvs.ConfigSource {
	p := s[KeyProps]
	if p == nil {
		panic("配置文件地址未初始化")
	}
	return p.(kvs.ConfigSource)
}

//资源启动器，每个应用少不了依赖其他资源，比如数据库，缓存，消息中间件等等服务
//启动器实现类，不需要实现所有方法，只需要实现对应的阶段方法即可，可以嵌入@BaseStarter
//通过实现资源启动器接口和资源启动注册器，友好的管理这些资源的初始化、安装、启动和停止。
//Starter对象注册器，所有需要在系统启动时需要实例化和运行的逻辑，都可以实现此接口
//注意只有Start方法才能被阻塞，如果是阻塞Start()，同时StartBlocking()要返回true
type Starter interface {
	// 系统启动 初始化一些基础资源
	Init(StarterContext)
	// 系统基础资源的安装
	Setup(StarterContext)
	// 启动基础资源
	Start(StarterContext)
	// 启动器是否可阻塞
	StartBlocking() bool
	// 资源停止和销毁
	Stop(StarterContext)
}

var _ Starter = new(BaseStarter)

//默认的空实现,方便资源启动器的实现
type BaseStarter struct {
}

func (b BaseStarter) Init(context StarterContext) {
}

func (b BaseStarter) Setup(context StarterContext) {
}

func (b BaseStarter) Start(context StarterContext) {
}

func (b BaseStarter) StartBlocking() bool {
	return false
}

func (b BaseStarter) Stop(context StarterContext) {
}

// 启动器注册器
type starterRegister struct {
	starters []Starter
}

// 启动器注册
func (r *starterRegister) Register(s Starter) {
	r.starters = append(r.starters, s)
}

func (r *starterRegister) AllStarters() []Starter {
	return r.starters
}

var StarterRegister *starterRegister = new(starterRegister)

func Register(s Starter) {
	StarterRegister.Register(s)
}
