# Stevedore
Simple tool that will add all your working clusters to you Codefresh account

# Install 
`go get github.com/codefresh-io/stevedore`

# Run 
Run `stevedore create --token {Codefresh JWT token}` to add all cluster from your `$HOME/.kube/config`

# Run as docker container
```
docker run \
-t \
-e KUBECONFIG=/config \
-e CODEFRESH_TOKEN=${PASTE_CODEFRESH_TOKEN} \
-e GOOGLE_APPLICATION_CREDENTIALS=/.config/gcloud/application_default_credentials.json \
-v ~/.kube/config:/config \
-v ~/.config/gcloud/:/.config/gcloud \
stevedore create
```
GOOGLE_APPLICATION_CREDENTIALS - needed when you have clusters that are hosted in Google

# Find you Codefresh JWT token
* Go to `https://g.codefresh.io/api/`
* Copy the token on the right side

## More functionallity
`stevedore create --token {Codefresh token} --config {another kube config valid file}`

# Todo:
* [ ] Tests!
* [ ] Support interactive mode
* [X] Support verbose/debug mode
* [ ] Support service-accounts from not default namespace
* [ ] Dry run


# Run as docker container
```
docker run \
-t \
-e KUBECONFIG=/config \
-e CODEFRESH_TOKEN=${PASTE_CODEFRESH_TOKEN} \
-e GOOGLE_APPLICATION_CREDENTIALS=/.config/gcloud/application_default_credentials.json \
-v ~/.kube/config:/config \
-v ~/.config/gcloud/:/.config/gcloud \
stevedore create
```
GOOGLE_APPLICATION_CREDENTIALS - needed when you have clusters that are hosted in Google
