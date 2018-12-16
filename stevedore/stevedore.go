package stevedore

import (
	"github.com/codefresh-io/stevedore/pkg/codefresh"
	"github.com/codefresh-io/stevedore/pkg/kubernetes"
	"github.com/codefresh-io/stevedore/pkg/reporter"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

func Init(c *cli.Context) {
	var name string
	codefreshAPI := codefresh.NewCodefreshAPI(c.String("api-host"), c.String("token"))
	reporter := reporter.NewReporter()
	kubernetesAPI := kubernetes.NewKubernetesAPI(c.String("config"), codefreshAPI, reporter)
	runOnAllContexts := c.IsSet("all")
	runOnContext := c.String("context")
	if c.IsSet("name-overwrite") {
		name = c.String("name-overwrite")
	} else {
		name = runOnContext
	}
	if runOnAllContexts {
		kubernetesAPI.GoOverAllContexts()
	} else if runOnContext != "" {
		kubernetesAPI.GoOverContextByName(runOnContext, c.String("namespace"), c.String("serviceaccount"), c.Bool("behind-firewall"), name)
	} else {
		kubernetesAPI.GoOverCurrentContext()
	}
	reporter.Print()
	log.Info("Operation is done, check your account setting")
}
