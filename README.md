# Stevedore
Simple tool that will add all your working clusters to you Codefresh account

# Install 
`go get github.com/codefresh-io/stevedore`

# Run 
Run `stevedore create --token {Codefresh JWT token}` to add all cluster from your `$HOME/.kube/config`

## More functionallity
`stevedore create --token {Codefresh token} --config {another kube config valid file}`


# Run in docker container
`docker run -v ~/.kube/config:/config codefresh/stevedore create --token {Codefresh token}  --config /config`

# Todo:
* [ ] Tests!
* [ ] Support interactive mode
* [ ] Support verbose/debug mode
* [ ] Support service-accounts from not default namespace