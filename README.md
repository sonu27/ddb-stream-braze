This lambda function takes in DynomoDB stream events and pushes the changes to Braze automatically.

### Environment Variables
- `BRAZE_API_KEY`

### Building
`GOOS=linux go build main.go`
`zip function.zip main`

Then upload the zip to AWS (will be automated later)
