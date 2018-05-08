package routers

import "github.com/astaxie/beego"
import "github.com/humpback/common/models"
import "github.com/humpback/humpback-agent/controllers"
import "github.com/humpback/humpback-agent/validators"

// Init - Init routers
func Init(composeStorage *models.ComposeStorage) {
	faqRouter := beego.NSRouter("/ping", &controllers.FaqController{})
	infoRouter := beego.NSRouter("/dockerinfo", &controllers.InfoController{})

	imageRouters := beego.NSNamespace("/images",
		beego.NSRouter("/", &controllers.ImageController{}, "get:GetImages;post:PullImage"),
		beego.NSRouter("/*", &controllers.ImageController{}, "get:GetImage;delete:DeleteImage"),
	)

	containerRouters := beego.NSNamespace("/containers",
		beego.NSRouter("/", &controllers.ContainerController{}, "get:GetContainers;post:CreateContainer;put:OperateContainer"),
		beego.NSRouter("/stats", &controllers.ContainerController{}, "get:GetAllContainerStats"),
		beego.NSRouter("/:containerid", &controllers.ContainerController{}, "get:GetContainer;delete:DeleteContainer"),
		beego.NSRouter("/:containerid/logs", &controllers.ContainerController{}, "get:GetContainerLogs"),
		beego.NSRouter("/:containerid/stats", &controllers.ContainerController{}, "get:GetContainerStats"),
		beego.NSRouter("/:containerid/status", &controllers.ContainerController{}, "get:GetContainerStatus"),
	)

	serviceRouters := beego.NSNamespace("/services",
		beego.NSRouter("/", &controllers.ServiceController{ComposeStorage: composeStorage}, "get:GetServices;post:CreateService;put:OperateService"),
		beego.NSRouter("/:service", &controllers.ServiceController{ComposeStorage: composeStorage}, "get:GetService;delete:DeleteService"),
		beego.NSRouter("/:service/upload", &controllers.ServiceController{ComposeStorage: composeStorage}, "post:PackageUploadService"),
	)

	ns := beego.NewNamespace("/dockerapi/v2",
		faqRouter,
		infoRouter,
		imageRouters,
		containerRouters,
		serviceRouters,
	)

	agentSpace := beego.NewNamespace("/v1",
		faqRouter,
		infoRouter,
		imageRouters,
		containerRouters,
		serviceRouters,
	)
	beego.AddNamespace(ns, agentSpace)

	beego.InsertFilter("/dockerapi/v2/containers", beego.BeforeExec, validators.CreateContainerValidator)
	beego.InsertFilter("/v1/containers", beego.BeforeExec, validators.CreateContainerValidator)
	beego.InsertFilter("/dockerapi/v2/services", beego.BeforeExec, validators.CreateServiceValidator)
	beego.InsertFilter("/v1/services", beego.BeforeExec, validators.CreateServiceValidator)
}
