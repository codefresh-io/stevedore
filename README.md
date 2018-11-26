# Stevedore
Simple tool that will add your working clusters to you Codefresh account

# Install 
`go get github.com/codefresh-io/stevedore`

# Run 
`stevedore create --help`
```
NAME:
   main create -

USAGE:
   main create [command options] [arguments...]

DESCRIPTION:
   Create clusters in Codefresh

OPTIONS:
   --verbose, -v              Turn on verbose mode
   --all, -a                  Add all clusters from config file, default is only current context
   --context value, -c value  Add spesific cluster
   --api-host value           Codefresh API host (default: "http://g.codefresh.io/") [$CODEFRESH_URL]
   --token value              Codefresh API Token [$CODEFRESH_TOKEN]
   --config value             Kubernetes config file to be used as input (default: "/Users/oleg/.kube/config") [$KUBECONFIG]
```

# Run as docker container
* No need to `go get github.com/codefresh-io/stevedore`
* Requiements:
    * Docker
    * kubeconfig placed in default directory `$HOME/.kube/config`
    * For clusters hosted in GKE - ensure `$HOME/.config/gcloud/application_default_credentials.json` exist
To add all availible clusters run:
```
docker run \
-t \
-e KUBECONFIG=/config \
-e CODEFRESH_TOKEN=$PASTE_CODEFRESH_TOKEN \
-e GOOGLE_APPLICATION_CREDENTIALS=/.config/gcloud/application_default_credentials.json \
-v ~/.kube/config:/config \
-v ~/.config/gcloud/:/.config/gcloud \
codefresh/stevedore create --all
```

# Generate Codefresh token
* https://g.codefresh.io/account-conf/tokens
