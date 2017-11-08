# Stevedore
Simple tool that will add all your working clusters to you Codefresh account

# Install 
`go get github.com/codefresh-io/stevedore`

# Run 
Run `stevedore create --token {Codefresh JWT token}` to add all cluster from your `$HOME/.kube/config`

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