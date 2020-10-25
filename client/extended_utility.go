package go_sqs_extended

import (
	"github.com/aws/aws-sdk-go/service/sqs"
	"sync"
	"sync/atomic"
)

func (esc *ExtendedSQS) messageIsLarge(message *sqs.SendMessageInput) bool {
	return int64(len([]byte(*message.MessageBody))) + esc.getMessageAttributesSize(message.MessageAttributes) >
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
