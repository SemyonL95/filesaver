## TASK
https://docs.google.com/document/d/1WJv0eUixKjRpWcV9Mbpfb3_UjtCyfR-AMdYVLoGUmvU/edit

## Deploy
1. Pull project
2. Run command "make app-server"
## Testing
1. Run command "make run-test"

After tests run command "app-down"
## Additional description
1. Project has two endpoints PUT "/upload" and GET "/files"
2. You could store multiple files which will store in storage folder with original name
2. "/files" endpoint required query string "filename", example: "/files?filename=test.txt", which will return file by name.
