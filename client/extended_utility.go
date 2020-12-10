package go_sqs_extended

import (
	"sync"
	"sync/atomic"

	"github.com/aws/aws-sdk-go/service/sqs"
)

func (esc *ExtendedSQS) messageIsLarge(message *sqs.SendMessageInput) bool {
	return int64(len([]byte(*message.MessageBody)))+esc.getMessageAttributesSize(message.MessageAttributes) >
		esc.cfg.LargeMessageThreshold
}

func (esc *ExtendedSQS) batchMessageIsLarge(message *sqs.SendMessageBatchInput) bool {
	sum := new(int64)
	var wg sync.WaitGroup
	for _, m := range message.Entries {
		wg.Add(1)
		go func(e *sqs.SendMessageBatchRequestEntry) {
			atomic.AddInt64(sum, int64(len([]byte(*e.MessageBody))))
			atomic.AddInt64(sum, esc.getMessageAttributesSize(e.MessageAttributes))
			wg.Done()
		}(m)
	}
	wg.Wait()
	return *sum > esc.cfg.LargeMessageThreshold
}

func (esc *ExtendedSQS) getMessageAttributesSize(attributes map[string]*sqs.MessageAttributeValue) int64 {
	sum := new(int64)
	var wg sync.WaitGroup
	for k, v := range attributes {
		wg.Add(1)
		go func(k string, attr *sqs.MessageAttributeValue) {
			atomic.AddInt64(sum, int64(len([]byte(k))))
			atomic.AddInt64(sum, int64(len(attr.BinaryValue)))
			wg.Done()
		}(k, v)
	}
	wg.Wait()
	return *sum
}

func (esc *ExtendedSQS) embedS3PointerInReceiptHandle(receiptHandle string, message *extendedQueueMessage) string {
	return S3BucketNameMarker + message.S3BucketName + S3BucketNameMarker +
		S3KeyMarker + message.S3Key + S3KeyMarker + receiptHandle
}

func (esc *ExtendedSQS) removeS3PointerFromReceiptHandle(receiptHandle string) string {
	index := ReceiptHandleRegexp.FindIndex([]byte(receiptHandle))
	if index == nil {
		return receiptHandle
	}
	return receiptHandle[index[1]:]
}

func (esc *ExtendedSQS) isS3PointerReceiptHandle(receiptHandle string) bool {
	return ReceiptHandleRegexp.Match([]byte(receiptHandle))
}

func (esc *ExtendedSQS) enforceSingleReserved(attrs []*string) []*string {
	var seen map[string]bool
	var attributes []*string
	for _, attr := range attrs {
		res := *attr
		if !ReservedAttribute(res).isReservedAttribute() {
			if seen[res] == false {
				attributes = append(attributes, &res)
				seen[res] = true
			}
		}
	}
	return attributes
}
