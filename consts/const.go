package consts

const ProjectField = "field.cattle.io/projectId"

var IgnoreUsers = []string{
	"kube-admin",
	"system:serviceaccount:cattle-system:kontainer-engine",
}
