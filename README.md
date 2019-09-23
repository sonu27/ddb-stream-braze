This lambda function takes in DynomoDB stream events and pushes the changes to Braze automatically.

### Environment Variables
- `BRAZE_API_KEY`

### Building
1. `GOOS=linux go build main.go`
2. `zip function.zip main`

Then upload the zip to AWS (will be automated later)
