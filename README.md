# go-sqs-extended
Extended Library for Amazon SQS (Simple Queue Service)

## What does this library do?

Amazon SQS has a message size limit of 256KiB per message. Sometimes, it's necessary to send messages that exceed this size limit.
This library is a wrapper over the Amazon SDK Client that will perform all the same operations, except...

- When function is enabled, messages exceeding a certain size limit will be uploaded to Amazon S3.
- The SQS message will contain a reference to the bucket and object like so: `{ "s3_bucket": "mybucket", "s3_object": "xyz.txt" }`
- When used as a consumer, this library will automatically retrieve large messages from S3.

The intent is a seamless development experience, where S3 is used as an intermediary buffer to send large payloads.
Small payloads will be sent "as-is", whereas larger ones will leverage S3. But to the client, the end result is the same!
