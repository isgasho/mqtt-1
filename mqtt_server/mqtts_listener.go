package main

import (
	"log"
	"mqtt"
)

type mqtts_listener struct {
	provider mqtt.Provider

	retainedMessages map[string]mqtt.Message
}

func newListener(provider mqtt.Provider) *mqtts_listener {
	return &mqtts_listener{provider: provider,
		retainedMessages: make(map[string]mqtt.Message)}
}
func (this *mqtts_listener) ProcessConnect(eventConnect mqtt.EventConnect) {
	log.Println("Received CONNECT")
	s := eventConnect.GetSession()

	pktconnack := mqtt.NewPacketConnack()
	pktconnack.SetSPFlag(false)
	pktconnack.SetReturnCode(mqtt.CONNACK_RETURNCODE_ACCEPTED)

	s.AcknowledgeConnect(pktconnack)
}
func (this *mqtts_listener) ProcessPublish(eventPublish mqtt.EventPublish) {
	log.Printf("Received Publish with DUP %v QoS %v RETAIN %v Topic: %v Content: %v\n",
		eventPublish.GetMessage().GetDup(),
		eventPublish.GetMessage().GetQos(),
		eventPublish.GetMessage().GetRetain(),
		eventPublish.GetMessage().GetTopic(),
		eventPublish.GetMessage().GetContent())

	this.provider.Forward(mqtt.NewMessage(eventPublish.GetMessage().GetDup(),
		eventPublish.GetMessage().GetQos(),
		false,
		eventPublish.GetMessage().GetTopic(),
		eventPublish.GetMessage().GetContent()))

	if eventPublish.GetMessage().GetRetain() {
		if eventPublish.GetMessage().GetContent() == "" {
			delete(this.retainedMessages, eventPublish.GetMessage().GetTopic())
		} else {
			this.retainedMessages[eventPublish.GetMessage().GetTopic()] = eventPublish.GetMessage()
		}
	}
}
func (this *mqtts_listener) ProcessSubscribe(eventSubscribe mqtt.EventSubscribe) {
	log.Printf("Received SUBSCRIBE with %v", eventSubscribe.GetSubscribeTopics())
	s := eventSubscribe.GetSession()

	pktsuback := mqtt.NewPacketSuback()
	pktsuback.SetPacketId(eventSubscribe.GetPacketId())
	topics := eventSubscribe.GetSubscribeTopics()
	qos := eventSubscribe.GetQoSs()
	retCodes := make([]byte, len(qos))
	for i := 0; i < len(qos); i++ {
		if topics[i] == "nosubscribe" {
			retCodes[i] = 0x80
		} else {
			retCodes[i] = byte(qos[i])
		}
	}
	pktsuback.SetReturnCodes(retCodes)

	s.AcknowledgeSubscribe(pktsuback)

	for _, msg := range this.retainedMessages {
		s.Forward(msg)
	}
}
func (this *mqtts_listener) ProcessUnsubscribe(eventUnsubscribe mqtt.EventUnsubscribe) {
	log.Printf("Received UNSUBSCRIBE with %v", eventUnsubscribe.GetUnsubscribeTopics())
}

func (this *mqtts_listener) ProcessTimeout(eventTimeout mqtt.EventTimeout) {
	log.Println("Received Timeout")
}
func (this *mqtts_listener) ProcessIOException(eventIOException mqtt.EventIOException) {
	log.Println("Received IOException")
}
func (this *mqtts_listener) ProcessSessionTerminated(eventSessionTerminated mqtt.EventSessionTerminated) {
	if eventSessionTerminated.GetWillMessage() != nil {
		log.Printf("Session Terminated with Reason: %s with WillMessage: DUP %v QoS %v RETAIN %v Topic: %v Content: %v\n",
			eventSessionTerminated.GetReason(),
			eventSessionTerminated.GetWillMessage().GetDup(),
			eventSessionTerminated.GetWillMessage().GetQos(),
			eventSessionTerminated.GetWillMessage().GetRetain(),
			eventSessionTerminated.GetWillMessage().GetTopic(),
			eventSessionTerminated.GetWillMessage().GetContent())

		this.provider.Forward(eventSessionTerminated.GetWillMessage())
	} else {
		log.Println("Session Terminated with Reason: ", eventSessionTerminated.GetReason(), "without WillMessage!")
	}
}
