# go-sqs-extended
Extended Library for Amazon SQS (Simple Queue Service)

## What does this library do?

Amazon SQS has a message size limit of 256KiB per message. A common use-case are sending dynamically sized messages thatmay  exceed this size limit.
This library is a wrapper over the Amazon SDK Client that will perform all the same operations, except...

- Messages exceeding a certain size limit will be uploaded to Amazon S3 instead.
- The SQS message will contain a reference to the bucket and object like so: `{ "s3_bucket": "mybucket", "s3_object": "xyz.txt" }`
- When used as a consumer, this library will automatically parse, retrieve, and return the contents of the S3 file.
- Messages under the size limit can be ignored and delivered as usual.

The intent is a seamless development experience, where S3 is used as an intermediary buffer to send large payloads.
