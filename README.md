# Video Transcoding Worker

Transcode video using `ffmpeg` written in Go

# Setting up

First, create a `config.yaml` file with following content:

```yaml
jwt_token: AdminJWTToken
endpoint: Backend endpoint
message_queue: rabbitmq_url
```

- `JWT Token`: The JWT token used to authenticate with the backend
- `Backend endpoint`: The endpoint of the backend server
- `RabbitMQ URL`: The URL of the RabbitMQ server

Then, run the following command to start the worker

```bash
go run main.go
```

# Important Notes

The default maximum retry count is 3. If the worker fails to transcode the video 3 times, it will stop retrying and send a message to the backend server to notify the failure.